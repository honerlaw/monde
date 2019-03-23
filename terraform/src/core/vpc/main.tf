resource "aws_internet_gateway" "core_internet_gateway" {
  vpc_id = "${aws_vpc.core_vpc.id}"

  tags {
    Name = "core internet gateway"
  }
}

resource "aws_vpc" "core_vpc" {
  cidr_block = "10.0.0.0/16"

  tags {
    Name = "core vpc"
  }
}

resource "aws_route_table" "core_public_route_table" {
  vpc_id = "${aws_vpc.core_vpc.id}"

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.core_internet_gateway.id}"
  }

  tags {
    Name = "Core Public Subnet Route Table"
  }
}

module "public_subnet_1" {
  source = "./subnet"
  vpc_id = "${aws_vpc.core_vpc.id}"
  cidr_prefix = "10.0.1"
  availability_zone = "us-east-1a"
  count = "1"
  route_table_id = "${aws_route_table.core_public_route_table.id}"
}

module "public_subnet_2" {
  source = "./subnet"
  vpc_id = "${aws_vpc.core_vpc.id}"
  cidr_prefix = "10.0.3"
  availability_zone = "us-east-1b"
  count = "2"
  route_table_id = "${aws_route_table.core_public_route_table.id}"
}

output "vpc_id" {
  value = "${aws_vpc.core_vpc.id}"
}

output "public_subnet_1_id" {
  value = "${module.public_subnet_1.public_subnet_id}"
}

output "public_subnet_2_id" {
  value = "${module.public_subnet_2.public_subnet_id}"
}

output "public_subnet_ids" {
  value = "${list(module.public_subnet_1.public_subnet_id, module.public_subnet_2.public_subnet_id)}"
}

output "public_subnet_cidr_blocks" {
  value = "${list(module.public_subnet_1.cidr_block, module.public_subnet_2.cidr_block)}"
}
