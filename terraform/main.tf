resource "aws_lambda_function" "lambda_function" {
  function_name = var.lambda_function_name
  filename      = "../lambda-handler.zip"
  handler       = "bootstrap"
  role          = aws_iam_role.lambda_role.arn
  source_code_hash = filebase64sha256("../lambda-handler.zip")
  runtime       = "provided.al2023"
  timeout       = 10
  environment {
    variables = {
      DISCORD_TOKEN = var.discord_token
    }
  }
}

resource "aws_scheduler_schedule" "start_poll_schedule" {
  name                         = var.start_poll_schedule_name
  schedule_expression          = var.start_poll_schedule_expression
  schedule_expression_timezone = var.time_zone
  flexible_time_window {
    mode                      = var.start_poll_schedule_flexible_time_window_mode
    maximum_window_in_minutes = var.start_poll_schedule_flexible_time_window_in_minutes
  }
  target {
    arn      = aws_lambda_function.lambda_function.arn
    role_arn = aws_iam_role.scheduler_role.arn
    input = jsonencode({
      action                = "startPoll"
      pollChannelId         = var.discord_poll_channel
      announcementChannelId = var.discord_announcement_channel
      timeZone              = var.time_zone
    })
  }
}

resource "aws_scheduler_schedule" "end_poll_schedule" {
  name                         = var.end_poll_schedule_name
  schedule_expression          = var.end_poll_schedule_expression
  schedule_expression_timezone = var.time_zone
  flexible_time_window {
    mode                      = var.end_poll_schedule_flexible_time_window_mode
    maximum_window_in_minutes = var.end_poll_schedule_flexible_time_window_in_minutes
  }
  target {
    arn      = aws_lambda_function.lambda_function.arn
    role_arn = aws_iam_role.scheduler_role.arn
    input = jsonencode({
      action                = "endPoll"
      pollChannelId         = var.discord_poll_channel
      announcementChannelId = var.discord_announcement_channel
      timeZone              = var.time_zone
    })
  }
}