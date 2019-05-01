variable "domain_zone" {}

resource "aws_route53_zone" "core_route53_zone" {
  name = "${var.domain_zone}"
}

output "route53_zone_id" {
  value = "${aws_route53_zone.core_route53_zone.id}"
}

output "route53_zone_name" {
  value = "${aws_route53_zone.core_route53_zone.name}"
}
