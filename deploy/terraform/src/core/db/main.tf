variable "vpc_id" {}
variable "subnet_cidr_blocks" {
  type = "list"
}
variable "subnet_ids" {
  type = "list"
}
variable "dbuser" {}
variable "dbpass" {}

resource "aws_security_group" "core_db_security_group" {
  vpc_id = "${var.vpc_id}"
  name = "core-db-security-group"
  description = "security group for the core aurora db access"

  ingress {
    protocol = "tcp"
    from_port = 3306
    to_port = 3306
    cidr_blocks = [
      "${var.subnet_cidr_blocks}"
    ]
  }

  egress {
    protocol = "-1"
    from_port = 0
    to_port = 0
    cidr_blocks = [
      "0.0.0.0/0"
    ]
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_db_subnet_group" "core_db_subnet_group" {
  name = "core-db-subnet-group"
  description = "core db subnet group"
  subnet_ids = [
    "${var.subnet_ids}"
  ]

  tags {
    Name = "core"
  }
}

resource "aws_rds_cluster" "core-rds-cluster" {
  cluster_identifier = "core-rds-cluster"
  vpc_security_group_ids = [
    "${aws_security_group.core_db_security_group.id}"
  ]
  db_subnet_group_name = "${aws_db_subnet_group.core_db_subnet_group.name}"
  engine_mode = "serverless"
  master_username = "${var.dbuser}"
  master_password = "${var.dbpass}"
  backup_retention_period = 7
  skip_final_snapshot = false

  scaling_configuration {
    auto_pause = true
    max_capacity = 2
    min_capacity = 2
    seconds_until_auto_pause = 300
  }

  lifecycle {
    ignore_changes = [
      "engine_vers¬ion"
    ]
  }
}