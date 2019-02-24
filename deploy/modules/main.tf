variable "region" {}
variable "image" {}

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

module "main_vpc" {
  source = "./vpc"
}

module "ecs_cluster" {
  source = "./ecs"
  vpc_id = "${module.main_vpc.vpc_id}"
  public_subnet_1_id = "${module.main_vpc.public_subnet_1_id}"
  image = "${var.image}"
}
