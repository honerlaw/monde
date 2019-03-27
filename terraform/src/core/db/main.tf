variable "vpc_id" {}
variable "subnet_cidr_blocks" {
  type = "list"
}
variable "subnet_ids" {
  type = "list"
}
variable "dbuser" {}
variable "dbpass" {}

resource "aws_secretsmanager_secret" "core_db_credential_secret" {
  name = "core-db-credential-secret"
  description = "core db credentials"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret_version" "core_db_credential_secret_version" {
  secret_id = "${aws_secretsmanager_secret.core_db_credential_secret.id}"
  secret_string = "{\"username\":\"${var.dbuser}\",\"password\":\"${var.dbpass}\"}"
}

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

  // @todo we really shouldn't do this
  ingress {
    protocol = "-1"
    from_port = 0
    to_port = 0
    cidr_blocks = [
      "0.0.0.0/0"
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

resource "aws_rds_cluster" "core_rds_cluster" {
  cluster_identifier = "core-rds-cluster"
  vpc_security_group_ids = [
    "${aws_security_group.core_db_security_group.id}"
  ]
  db_subnet_group_name = "${aws_db_subnet_group.core_db_subnet_group.name}"
  engine_mode = "provisioned"
  master_username = "${var.dbuser}"
  master_password = "${var.dbpass}"
  backup_retention_period = 7
  skip_final_snapshot = false
  final_snapshot_identifier = "core-rds-final-snapshot-${replace(timestamp(), ":", "-")}"

  lifecycle {
    create_before_destroy = true
    ignore_changes = [
      "engine_version"
    ]
  }
}

resource "aws_rds_cluster_instance" "core_rds_cluster_instance" {
  count = 1
  identifier = "core-rds-instance-1"
  cluster_identifier = "${aws_rds_cluster.core_rds_cluster.id}"
  instance_class = "db.t2.small"
  db_subnet_group_name = "${aws_db_subnet_group.core_db_subnet_group.name}"
  publicly_accessible = true

  lifecycle {
    create_before_destroy = true
  }
}

output "core_rds_cluster_endpoint" {
  value = "${aws_rds_cluster.core_rds_cluster.endpoint}"
}

