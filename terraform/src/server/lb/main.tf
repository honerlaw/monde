variable subnets {
  type = "list"
}
variable vpc_id {}
variable container_port {}

locals {
  bucket_name = "vueon-server-lb-access-logs"
  log_prefix = "logs"
}

resource "aws_lb_target_group" "server_lb_target_group" {
  target_type = "ip"
  name = "server-lb-target-group"
  port = "${var.container_port}"
  protocol = "HTTP"
  vpc_id = "${var.vpc_id}"
  depends_on = ["aws_lb.server_lb"]
  health_check {
    path = "/health"
  }
}

data "aws_iam_policy_document" "server_lb_policy_document" {
  policy_id = "s3_lb_write"

  statement = {
    actions = ["s3:PutObject"]
    resources = ["arn:aws:s3:::${local.bucket_name}/${local.log_prefix}/*"]

    principals = {
      identifiers = ["*"]
      type = "AWS"
    }
  }
}

resource "aws_s3_bucket" "server_lb_log_bucket" {
  bucket = "${local.bucket_name}"
  policy = "${data.aws_iam_policy_document.server_lb_policy_document.json}"
}

resource "aws_security_group" "server_lb_security_group" {
  name = "server-lb-security-group"
  vpc_id = "${var.vpc_id}"
  ingress {
    from_port = "${var.container_port}"
    protocol = "tcp"
    to_port = "${var.container_port}"
    cidr_blocks = [
      "0.0.0.0/0"
    ]
  }
  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = [
      "0.0.0.0/0"
    ]
  }
}

resource "aws_lb" "server_lb" {
  name = "server-lb"
  internal = false
  load_balancer_type = "application"
  subnets = ["${var.subnets}"]
  access_logs {
    prefix = "${local.log_prefix}"
    bucket = "${aws_s3_bucket.server_lb_log_bucket.bucket}"
    enabled = true
  }
  security_groups = ["${aws_security_group.server_lb_security_group.id}"]
}

resource "aws_lb_listener" "https_listener" {
  load_balancer_arn = "${aws_lb.server_lb.arn}"
  port = "80"
  protocol = "HTTP"

  default_action {
    type = "forward"
    target_group_arn = "${aws_lb_target_group.server_lb_target_group.arn}"
  }
}

output "lb_target_group_arn" {
  value = "${aws_lb_target_group.server_lb_target_group.arn}"
}
