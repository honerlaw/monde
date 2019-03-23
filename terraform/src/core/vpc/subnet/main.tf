variable "vpc_id" {}
variable "cidr_prefix" {}
variable "count" {}
variable "route_table_id" {}
variable "availability_zone" {}

locals {
  cidr_block = "${var.cidr_prefix}.0/24"
}

resource "aws_subnet" "core_public_subnet" {
  vpc_id = "${var.vpc_id}"
  cidr_block = "${local.cidr_block}"
  availability_zone = "${var.availability_zone}"
  map_public_ip_on_launch = true

  tags {
    Name = "core public subnet ${var.count}"
  }
}

resource "aws_route_table_association" "core_public_route_table_association" {
  subnet_id = "${aws_subnet.core_public_subnet.id}"
  route_table_id = "${var.route_table_id}"
}

output "public_subnet_id" {
  value = "${aws_subnet.core_public_subnet.id}"
}

output "cidr_block" {
  value = "${local.cidr_block}"
}