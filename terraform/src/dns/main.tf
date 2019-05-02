variable "region" {}
variable "domain_zone" {}
variable "email_domain" {}

terraform {
  backend "s3" {
    bucket = "vueon-terraform"
    key = "terraform/dns.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = "${var.region}"
}

module "core_route53" {
  source = "./route53"
  domain_zone = "${var.domain_zone}"
}

module "core_ses" {
  source = "./ses"
  region = "${var.region}"
  zone_id = "${module.core_route53.route53_zone_id}"
  domain = "${var.email_domain}"
}
