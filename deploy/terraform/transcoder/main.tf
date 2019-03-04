variable "region" {}

terraform {
  backend "s3" {
    bucket = "bddevelop-monde"
    key = "terraform/transcoder.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = "${var.region}"
}

resource "aws_s3_bucket" "video_upload" {
  bucket = "monde-aws-video-upload"
}

resource "aws_s3_bucket" "video_processed" {
  bucket = "monde-aws-video-processed"
}

resource "aws_s3_bucket" "video_thumbnails" {
  bucket = "monde-aws-video-thumbnails"
}

resource "aws_elastictranscoder_pipeline" "video_pipeline" {
  input_bucket = "${aws_s3_bucket.video_upload.bucket}"
  name = "aws_elastictranscoder_video_pipeline"
  role = "${aws_iam_role.video_pipeline_role.arn}"

  content_config {
    bucket = "${aws_s3_bucket.video_processed.bucket}"
    storage_class = "Standard"
  }

  thumbnail_config {
    bucket = "${aws_s3_bucket.video_thumbnails.bucket}"
    storage_class = "Standard"
  },

  depends_on = [
    "aws_s3_bucket.video_upload",
    "aws_s3_bucket.video_processed",
    "aws_s3_bucket.video_thumbnails"]
}

data "aws_iam_policy_document" "video_pipeline_policy_document" {
  statement {
    effect = "Allow"
    actions = [
      "s3:Get*",
      "s3:ListBucket",
      "s3:Put*",
      "s3:*MultipartUpload*"
    ]
    resources = ["*"]
  }
  statement {
    effect = "Allow"
    actions = [
      "sns:Publish"
    ]
    resources = ["*"]
  }
  statement {
    effect = "Deny"
    actions = [
      "sns:*Permission*",
      "sns:*Delete*",
      "sns:*Remove*",
      "s3:*Policy*",
      "s3:*Delete*"
    ]
    resources = ["*"]
  }
}

resource "aws_iam_policy" "video_pipeline_role_policy" {
  name = "video-pipeline-role-policy"
  path = "/"
  policy = "${data.aws_iam_policy_document.video_pipeline_policy_document.json}"
}

resource "aws_iam_role_policy_attachment" "video_pipeline_role_policy_attachment" {
  role = "${aws_iam_role.video_pipeline_role.name}"
  policy_arn = "${aws_iam_policy.video_pipeline_role_policy.arn}"
}

resource "aws_iam_role" "video_pipeline_role" {
  name = "ecs_execution_role"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "ec2.amazonaws.com"
        ]
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}