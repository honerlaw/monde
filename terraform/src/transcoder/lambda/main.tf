variable "upload_bucket_id" {}
variable "upload_bucket_arn" {}
variable "lambda_file_path" {}

resource "aws_iam_role" "transcoder_lambda_iam_role" {
  name = "transcoder-lambda-iam-role"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "transcoder_lambda_policy" {
  name = "transcoder-lambda-policy"
  path = "/"
  description = "policy for logging from the transcoder lambda"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "*"
      ],
      "Resource": "*",
      "Effect": "Allow"
    },
    {
      "Action": [
        "logs:CreateLogStream",
        "logs:PutLogEvents",
        "logs:CreateLogGroup"
      ],
      "Resource": "arn:aws:logs:*:*:*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "transcoder_lambda_policy_attachment" {
  role = "${aws_iam_role.transcoder_lambda_iam_role.name}"
  policy_arn = "${aws_iam_policy.transcoder_lambda_policy.arn}"
}

resource "aws_lambda_function" "transcoder_lambda_function" {
  filename = "${var.lambda_file_path}"
  function_name = "transcoder-lambda"
  role = "${aws_iam_role.transcoder_lambda_iam_role.arn}"
  handler = "bin/transcoder"
  source_code_hash = "${base64sha256(file("${var.lambda_file_path}"))}"
  runtime = "go1.x"
}

resource "aws_lambda_permission" "transcoder_lambda_allow_upload_bucket" {
  statement_id = "AllowExecutionFromS3Bucket"
  action = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.transcoder_lambda_function.arn}"
  principal = "s3.amazonaws.com"
  source_arn = "${var.upload_bucket_arn}"
}

resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = "${var.upload_bucket_id}"
  lambda_function {
    lambda_function_arn = "${aws_lambda_function.transcoder_lambda_function.arn}"
    events = [
      "s3:ObjectCreated:*"
    ]
  }
}