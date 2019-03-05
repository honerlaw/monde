variable "region" {}

terraform {
  backend "s3" {
    bucket = "monde-terraform"
    key = "terraform/transcoder.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = "${var.region}"
}

module "transcoder_iam" {
  source = "./iam"
}

resource "aws_s3_bucket" "transcoder_upload" {
  bucket = "monde-aws-transcoder-upload"
}

resource "aws_s3_bucket" "transcoder_processed" {
  bucket = "monde-aws-transcoder-processed"
}

resource "aws_s3_bucket" "transcoder_thumbnails" {
  bucket = "monde-aws-transcoder-thumbnails"
}

resource "aws_elastictranscoder_pipeline" "transcoder_pipeline" {
  input_bucket = "${aws_s3_bucket.transcoder_upload.bucket}"
  name = "transcoder-pipeline"
  role = "${module.transcoder_iam.transcoder_pipeline_role_arn}"

  content_config {
    bucket = "${aws_s3_bucket.transcoder_processed.bucket}"
    storage_class = "Standard"
  }

  thumbnail_config {
    bucket = "${aws_s3_bucket.transcoder_thumbnails.bucket}"
    storage_class = "Standard"
  },

  depends_on = [
    "aws_s3_bucket.transcoder_upload",
    "aws_s3_bucket.transcoder_processed",
    "aws_s3_bucket.transcoder_thumbnails"
  ]
}