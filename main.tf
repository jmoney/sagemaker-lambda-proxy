terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "5.16.1"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

variable "endpoint_name" {
  type = string
  description = "value of the endpoint name to the sagemaker endpoint"
}

data "aws_region" "current" {}

resource "aws_iam_policy" "sagemaker_lambda_proxy" {
  name = "SagemakerLambdaProxy${data.aws_region.current.name}"
  description = "Policy for Sagemaker Lambda Proxy"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = [
          "sagemaker:InvokeEndpoint"
        ],
        Effect = "Allow",
        Resource = "*"
      },
      {
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ],
        Effect = "Allow",
        Resource = "arn:aws:logs:*:*:*"
      }
    ]
  })
}

resource "aws_iam_role" "lambda" {
  name = "SagemakerLambdaProxy${data.aws_region.current.name}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole",
        Principal = {
          Service = "lambda.amazonaws.com"
        },
        Effect = "Allow",
        Sid = "AllowLambda"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_role_attachment" {
  role = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.sagemaker_lambda_proxy.arn
}

data "archive_file" "archive" {
  type        = "zip"
  source_file = "${path.root}/bin/bootstrap"
  output_file_mode = "0666"
  output_path = "${path.module}/sagemaker-lambda-proxy.zip"
}

resource "random_uuid" "nonce" {
  keepers = {
    version = 1
  }
}

resource "aws_lambda_function" "sagemaker_lambda_proxy" {
  function_name = "SagemakerLambdaProxy"
  description = "Lambda that proxies the request to SageMaker"
  role = aws_iam_role.lambda.arn
  handler = "bootstrap"
  runtime = "provided.al2"
  package_type = "Zip"
  timeout = 900
  filename = data.archive_file.archive.output_path
  source_code_hash = data.archive_file.archive.output_base64sha256

  environment {
    variables = {
      ENDPOINT_NAME = var.endpoint_name
      NONCE = random_uuid.nonce.result
    }
  }
}

resource "aws_apigatewayv2_api" "sagemaker_proxy_api" {
    name          = "api-sagemaker-proxy"
    protocol_type = "HTTP"
    target        = aws_lambda_function.sagemaker_lambda_proxy.arn
}

resource "aws_lambda_permission" "apigw" {
    action        = "lambda:InvokeFunction"
    function_name = aws_lambda_function.sagemaker_lambda_proxy.arn
    principal     = "apigateway.amazonaws.com"
    source_arn = "${aws_apigatewayv2_api.sagemaker_proxy_api.execution_arn}/*/*"
}

output "nonce" {
    sensitive = true
    value = random_uuid.nonce.result
}

output "url" {
    value = aws_apigatewayv2_api.sagemaker_proxy_api.api_endpoint
}