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

resource "aws_subnet" "core_public_subnet_1" {
  vpc_id = "${aws_vpc.core_vpc.id}"
  cidr_block = "10.0.1.0/24"
  availability_zone = "us-east-1a"
  map_public_ip_on_launch = true
  depends_on = ["aws_internet_gateway.core_internet_gateway"]

  tags {
    Name = "core public subnet 1"
  }
}

resource "aws_subnet" "core_public_subnet_2" {
  vpc_id = "${aws_vpc.core_vpc.id}"
  availability_zone = "us-east-1b"
  cidr_block = "10.0.2.0/24"
  map_public_ip_on_launch = true
  depends_on = ["aws_internet_gateway.core_internet_gateway"]

  tags {
    Name = "core public subnet 2"
  }
}

output "vpc_id" {
  value = "${aws_vpc.core_vpc.id}"
}

output "public_subnet_1_id" {
  value = "${aws_subnet.core_public_subnet_1.id}"
}

output "public_subnet_2_id" {
  value = "${aws_subnet.core_public_subnet_2.id}"
}

output "public_subnet_ids" {
  value = "${list(aws_subnet.core_public_subnet_1.id, aws_subnet.core_public_subnet_2.id)}"
}

output "public_subnet_cidr_blocks" {
  value = "${list(aws_subnet.core_public_subnet_1.cidr_block, aws_subnet.core_public_subnet_2.cidr_block)}"
}
