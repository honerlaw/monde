variable subnets {
  type = "list"
}
variable vpc_id {}
variable container_port {}

resource "aws_lb_target_group" "server_lb_target_group" {
  target_type = "ip"
  name = "server-lb-target-group"
  port = "${var.container_port}"
  protocol = "HTTP"
  vpc_id = "${var.vpc_id}"
  depends_on = ["aws_lb.server_lb"]
}

resource "aws_lb" "server_lb" {
  name = "server-lb"
  internal = false
  load_balancer_type = "application"
  subnets = ["${var.subnets}"]
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
