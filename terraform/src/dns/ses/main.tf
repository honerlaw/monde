variable "region" {}
variable "domain" {}
variable "zone_id" {}

resource "aws_ses_domain_identity" "domain" {
  domain = "${var.domain}"
}

resource "aws_ses_domain_mail_from" "domain_mail_from" {
  domain           = "${aws_ses_domain_identity.domain.domain}"
  mail_from_domain = "bounce.${aws_ses_domain_identity.domain.domain}"
}

resource "aws_ses_domain_dkim" "dkim" {
  domain = "${aws_ses_domain_identity.domain.domain}"
}

resource "aws_route53_record" "domain_amazonses_verification_record" {
  count   = "1"
  zone_id = "${var.zone_id}"
  name    = "_amazonses.${var.domain}"
  type    = "TXT"
  ttl     = "3600"
  records = ["${aws_ses_domain_identity.domain.verification_token}"]
}

resource "aws_route53_record" "domain_amazonses_dkim_record" {
  count   = "3"
  zone_id = "${var.zone_id}"
  name    = "${element(aws_ses_domain_dkim.dkim.dkim_tokens, count.index)}._domainkey.${var.domain}"
  type    = "CNAME"
  ttl     = "3600"
  records = ["${element(aws_ses_domain_dkim.dkim.dkim_tokens, count.index)}.dkim.amazonses.com"]
}

output "domain_identity_arn" {
  description = "ARN of the SES domain identity"
  value       = "${aws_ses_domain_identity.domain.arn}"
}

resource "aws_route53_record" "ses_domain_mail_from_mx" {
  zone_id = "${var.zone_id}"
  name    = "${aws_ses_domain_mail_from.domain_mail_from.mail_from_domain}"
  type    = "MX"
  ttl     = "600"
  records = ["10 feedback-smtp.${var.region}.amazonses.com"]
}

resource "aws_route53_record" "ses_domain_mail_from_txt" {
  zone_id = "${var.zone_id}"
  name    = "${aws_ses_domain_mail_from.domain_mail_from.mail_from_domain}"
  type    = "TXT"
  ttl     = "600"
  records = ["v=spf1 include:amazonses.com -all"]
}
