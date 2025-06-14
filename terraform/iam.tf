resource "aws_iam_role" "lambda_role" {
  name               = "lambda-role"
  assume_role_policy = data.aws_iam_policy_document.assume_role_lambda.json
}

data "aws_iam_policy_document" "assume_role_lambda" {
  statement {
    effect = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      type = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "scheduler_role" {
  name               = "eventbridge-scheduler-role"
  assume_role_policy = data.aws_iam_policy_document.assume_role_scheduler.json
}

data "aws_iam_policy_document" "assume_role_scheduler" {
  statement {
    effect = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      type = "Service"
      identifiers = ["scheduler.amazonaws.com"]
    }
  }
}

resource "aws_iam_role_policy" "scheduler_lambda_policy" {
  name   = "eventbridge-scheduler-lambda-policy"
  role   = aws_iam_role.scheduler_role.id
  policy = data.aws_iam_policy_document.invoke_function_lambda.json
}

data "aws_iam_policy_document" "invoke_function_lambda" {
  statement {
    effect = "Allow"
    actions = ["lambda:InvokeFunction"]
    resources = [aws_lambda_function.lambda_function.arn]
  }
}