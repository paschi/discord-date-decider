resource "aws_lambda_function" "go_function" {
  function_name = "discord-date-decider"
  filename      = "../lambda-handler.zip"
  handler       = "bootstrap"
  role          = aws_iam_role.iam_for_lambda.arn
  source_code_hash = filebase64sha256("../lambda-handler.zip")
  runtime       = "provided.al2023"
  architectures = ["x86_64"]
  timeout       = 10
  environment {
    variables = {
      DISCORD_TOKEN = var.discord_token
    }
  }
}

resource "aws_lambda_invocation" "test_event" {
  function_name = aws_lambda_function.go_function.function_name
  input = jsonencode({
    action              = "startPoll"
    pollChannel         = var.discord_channel
    announcementChannel = var.discord_channel
  })
  triggers = {
    redeployment = aws_lambda_function.go_function.source_code_hash
  }
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"
    principals {
      type = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "iam_for_lambda" {
  name               = "iam_for_lambda"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}
