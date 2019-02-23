variable "image" {}
variable "vpc_id" {}
variable "public_subnet_1_id" {}

resource "aws_ecs_cluster" "monde_ecs_cluster" {
  name = "monde-ecs-cluster"
}

resource "aws_security_group" "monde_ecs_security_group" {
  vpc_id = "${var.vpc_id}"
  ingress {
    from_port = 80
    protocol = "http"
    to_port = 80
  }
}

resource "aws_ecs_service" "monde_ecs_service" {
  name = "monde-ecs-server"
  cluster = "${aws_ecs_cluster.monde_ecs_cluster.id}"
  desired_count = 1
  launch_type = "FARGATE"
  network_configuration {
    security_groups = [
      "${aws_security_group.monde_ecs_security_group.id}"
    ]
    subnets = [
      "${var.public_subnet_1_id}"
    ]
  }
  task_definition = "${aws_ecs_task_definition.monde_ecs_task_definition.family}:${max("${aws_ecs_task_definition.monde_ecs_task_definition.revision}", "${aws_ecs_task_definition.monde_ecs_task_definition.revision}")}"
}

resource "aws_ecs_task_definition" "monde_ecs_task_definition" {
  family = "monnde-ecs-family"
  network_mode = "awsvpc"
  cpu = "512"
  memory = "1024"
  requires_compatibilities = [
    "FARGATE"
  ]
  container_definitions = "${data.template_file.monde_ecs_container_definitions.rendered}"
}

data "template_file" "monde_ecs_container_definitions" {
  template = "${file("${path.module}/container_definitions.json")}"
  vars {
    image = "${var.image}"
  }
}