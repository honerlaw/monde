resource "aws_vpc" "main_vpc" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "main_public_subnet_1" {
  vpc_id = "${aws_vpc.main_vpc.id}"
  cidr_block = "10.0.1.0/24"
  map_public_ip_on_launch = true
}

output "vpc_id" {
  value = "${aws_vpc.main_vpc.id}"
}

output "public_subnet_1_id" {
  value = "${aws_subnet.main_public_subnet_1.id}"
}