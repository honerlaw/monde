variable "region" {}
variable "image" {}
variable "container_name"{}
variable "container_port"{}

terraform {
  backend "s3" {
    bucket = "vueon-terraform"
    key = "terraform/server.tfstate"
    region = "us-east-1"
  }
}

data "terraform_remote_state" "core" {
  backend = "s3"
  config {
    bucket = "vueon-terraform"
    key = "terraform/core.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = "${var.region}"
}

module "server_lb" {
  source = "./lb"
  subnets = "${data.terraform_remote_state.core.public_subnet_ids}"
  vpc_id = "${data.terraform_remote_state.core.vpc_id}"
  container_port = "${var.container_port}"
}

module "server_cluster" {
  source = "./ecs"
  vpc_id = "${data.terraform_remote_state.core.vpc_id}"
  public_subnet_1_id = "${data.terraform_remote_state.core.public_subnet_1_id}"
  image = "${var.image}"
  lb_target_group_arn = "${module.server_lb.lb_target_group_arn}"
  container_name = "${var.container_name}"
  container_port = "${var.container_port}"
  region = "${var.region}"
}
