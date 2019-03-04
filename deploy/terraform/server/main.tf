variable "region" {}
variable "image" {}
variable "container_name"{}
variable "container_port"{}

terraform {
  backend "s3" {
    bucket = "bddevelop-monde"
    key = "terraform/server.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = "${var.region}"
}

module "main_vpc" {
  source = "./vpc"
}

module "main_lb" {
  source = "./lb"
  subnets = "${module.main_vpc.public_subnet_ids}"
  vpc_id = "${module.main_vpc.vpc_id}"
  container_port = "${var.container_port}"
}

module "ecs_cluster" {
  source = "./ecs"
  vpc_id = "${module.main_vpc.vpc_id}"
  public_subnet_1_id = "${module.main_vpc.public_subnet_1_id}"
  image = "${var.image}"
  lb_target_group_arn = "${module.main_lb.lb_target_group_arn}"
  container_name = "${var.container_name}"
  container_port = "${var.container_port}"
}
