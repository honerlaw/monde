variable "region" {}
variable "dbuser" {}
variable "dbpass" {}

terraform {
  backend "s3" {
    bucket = "vueon-terraform"
    key = "terraform/core.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = "${var.region}"
}

module "core_vpc" {
  source = "./vpc"
}

module "core_db" {
  source = "./db"
  vpc_id = "${module.core_vpc.vpc_id}"
  subnet_cidr_blocks = "${module.core_vpc.public_subnet_cidr_blocks}"
  subnet_ids = "${module.core_vpc.public_subnet_ids}"
  dbuser = "${var.dbuser}"
  dbpass = "${var.dbpass}"
}

output "vpc_id" {
  value = "${module.core_vpc.vpc_id}"
}

output "public_subnet_ids" {
  value = "${module.core_vpc.public_subnet_ids}"
}

output "public_subnet_cidr_blocks" {
  value = "${module.core_vpc.public_subnet_cidr_blocks}"
}

output "public_subnet_1_id" {
  value = "${module.core_vpc.public_subnet_1_id}"
}

output "public_subnet_2_id" {
  value = "${module.core_vpc.public_subnet_2_id}"
}

output "rds_cluster_endpoint" {
  value = "${module.core_db.core_rds_cluster_endpoint}"
}
