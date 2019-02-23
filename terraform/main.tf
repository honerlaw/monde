terraform {
  backend "s3" {
    bucket = "bddevelop-monde"
    key = "terraform/terraform.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = "${var.region}"
}

resource "aws_vpc" "main_vpc" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "main_public_subnet_1" {
  vpc_id = "${aws_vpc.main_vpc.id}"
  cidr_block = "10.0.1.0/24"
  map_public_ip_on_launch = true
}

resource "aws_security_group" "main_security_group" {
  vpc_id = "${aws_vpc.main_vpc.id}"
}

resource "aws_ecs_cluster" "monde_ecs_cluster" {
  name = "monde-ecs-cluster"
}

resource "aws_ecs_service" "monde_ecs_service" {
  name = "monde-ecs-server"
  cluster = "${aws_ecs_cluster.monde_ecs_cluster.id}"
  desired_count = 1
  launch_type = "FARGATE"
  network_configuration {
    security_groups = [
      "${aws_security_group.main_security_group.id}"
    ]
    subnets = [
      "${aws_subnet.main_public_subnet_1.id}"
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
  container_definitions = "${data.template_file.container_definitions.rendered}"
}

data "template_file" "container_definitions" {
  template = "${file("${path.module}/container_definitions.json")}"
  vars {
    image = "${var.image}"
  }
}
