data "aws_iam_policy_document" "transcoder_pipeline_policy_document" {
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

resource "aws_iam_policy" "transcoder_pipeline_role_policy" {
  name = "transcoder-pipeline-role-policy"
  path = "/"
  policy = "${data.aws_iam_policy_document.transcoder_pipeline_policy_document.json}"
}

resource "aws_iam_role_policy_attachment" "transcoder_pipeline_role_policy_attachment" {
  role = "${aws_iam_role.transcoder_pipeline_role.name}"
  policy_arn = "${aws_iam_policy.transcoder_pipeline_role_policy.arn}"
}

resource "aws_iam_role" "transcoder_pipeline_role" {
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

output "transcoder_pipeline_role_arn" {
  value = "${aws_iam_role.transcoder_pipeline_role.arn}"
}
