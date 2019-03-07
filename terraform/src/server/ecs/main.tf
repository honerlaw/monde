variable "image" {}
variable "vpc_id" {}
variable "public_subnet_1_id" {}
variable "container_name"{}
variable "container_port"{}
variable "lb_target_group_arn"{}

resource "aws_ecs_cluster" "ecs_cluster" {
  name = "ecs_cluster"
}

resource "aws_security_group" "ecs_security_group" {
  vpc_id = "${var.vpc_id}"
  ingress {
    from_port = 80
    protocol = "http"
    to_port = 80
  }
}

resource "aws_ecs_service" "ecs_service" {
  name = "ecs_server"
  cluster = "${aws_ecs_cluster.ecs_cluster.id}"
  desired_count = 1
  launch_type = "FARGATE"
  load_balancer {
    target_group_arn = "${var.lb_target_group_arn}"
    container_name = "${var.container_name}"
    container_port = "${var.container_port}"
  }
  network_configuration {
    security_groups = [
      "${aws_security_group.ecs_security_group.id}"
    ]
    subnets = [
      "${var.public_subnet_1_id}"
    ]
  }
  task_definition = "${aws_ecs_task_definition.ecs_task_definition.family}:${max("${aws_ecs_task_definition.ecs_task_definition.revision}", "${aws_ecs_task_definition.ecs_task_definition.revision}")}"
}

resource "aws_ecs_task_definition" "ecs_task_definition" {
  family = "ecs_family"
  network_mode = "awsvpc"
  cpu = "512"
  memory = "1024"
  execution_role_arn = "${aws_iam_role.ecs_execution_role.arn}"
  requires_compatibilities = [
    "FARGATE"
  ]
  container_definitions = "${data.template_file.ecs_container_definitions.rendered}"
}

data "template_file" "ecs_container_definitions" {
  template = "${file("${path.module}/container_definitions.json")}"
  vars {
    image = "${var.image}"
  }
}

data "aws_iam_policy_document" "ecs_execution_role_policy_document" {
  statement {
    actions = [
      "*",
    ]
    resources = [
      "*"
    ]
  }
}

data "aws_iam_policy_document" "ecs_execution_assume_role_policy_document" {
  statement {

    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }
  }
}

resource "aws_iam_policy" "ecs_execution_role_policy" {
  name   = "ecs_execution_role_policy"
  path   = "/"
  policy = "${data.aws_iam_policy_document.ecs_execution_role_policy_document.json}"
}

resource "aws_iam_role_policy_attachment" "ecs_execution_role_policy_attachment" {
  role = "${aws_iam_role.ecs_execution_role.name}"
  policy_arn = "${aws_iam_policy.ecs_execution_role_policy.arn}"
}

resource "aws_iam_role" "ecs_execution_role" {
  name = "ecs_execution_role"
  assume_role_policy = "${data.aws_iam_policy_document.ecs_execution_assume_role_policy_document.json}"
}