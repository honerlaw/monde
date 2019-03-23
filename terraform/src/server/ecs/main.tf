variable "image" {}
variable "vpc_id" {}
variable "public_subnet_1_id" {}
variable "container_name" {}
variable "container_port" {}
variable "lb_target_group_arn" {}
variable "region" {}

data "aws_iam_policy_document" "server_execution_role_policy_document" {
  statement {
    actions = [
      "*",
    ]
    resources = [
      "*"
    ]
  }
}

data "aws_iam_policy_document" "server_execution_assume_role_policy_document" {
  statement {

    actions = [
      "sts:AssumeRole"
    ]

    principals {
      type = "Service"
      identifiers = [
        "ec2.amazonaws.com"
      ]
    }

    principals {
      type = "AWS"
      identifiers = [
        "*"
      ]
    }
  }
}

resource "aws_iam_policy" "server_execution_role_policy" {
  name = "server-execution-role-policy"
  path = "/"
  policy = "${data.aws_iam_policy_document.server_execution_role_policy_document.json}"
}

resource "aws_iam_role_policy_attachment" "server_execution_role_policy_attachment" {
  role = "${aws_iam_role.server_execution_role.name}"
  policy_arn = "${aws_iam_policy.server_execution_role_policy.arn}"
}

resource "aws_iam_role" "server_execution_role" {
  name = "server-execution-role"
  assume_role_policy = "${data.aws_iam_policy_document.server_execution_assume_role_policy_document.json}"
}

resource "aws_security_group" "server_security_group" {
  name = "server-security-group"
  vpc_id = "${var.vpc_id}"
  ingress {
    from_port = "${var.container_port}"
    protocol = "tcp"
    to_port = "${var.container_port}"
    cidr_blocks = [
      "0.0.0.0/0"
    ]
  }
  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = [
      "0.0.0.0/0"
    ]
  }
}

resource "aws_cloudwatch_log_group" "server_log_group" {
  name = "awslogs-vueon"
}

resource "aws_ecs_cluster" "server_cluster" {
  name = "server_cluster"
}

data "template_file" "server_container_definitions" {
  template = "${file("${path.module}/container_definitions.json")}"
  vars {
    image = "${var.image}"
    container_port = "${var.container_port}"
    region = "${var.region}"
    log_group = "${aws_cloudwatch_log_group.server_log_group.name}"
  }
}

resource "aws_ecs_task_definition" "server_task_definition" {
  family = "server-family"
  network_mode = "awsvpc"
  cpu = "512"
  memory = "1024"
  execution_role_arn = "${aws_iam_role.server_execution_role.arn}"
  requires_compatibilities = [
    "FARGATE"
  ]
  depends_on = ["aws_cloudwatch_log_group.server_log_group"]
  container_definitions = "${data.template_file.server_container_definitions.rendered}"
}

resource "aws_ecs_service" "server_service" {
  name = "server-service"
  cluster = "${aws_ecs_cluster.server_cluster.id}"
  desired_count = 1
  launch_type = "FARGATE"
  load_balancer {
    target_group_arn = "${var.lb_target_group_arn}"
    container_name = "${var.container_name}"
    container_port = "${var.container_port}"
  }
  network_configuration {
    assign_public_ip = true
    security_groups = [
      "${aws_security_group.server_security_group.id}"
    ]
    subnets = [
      "${var.public_subnet_1_id}"
    ]
  }
  task_definition = "${aws_ecs_task_definition.server_task_definition.arn}"
}
