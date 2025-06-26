variable "discord_announcement_channel" {
  type = string
}

variable "discord_poll_channel" {
  type = string
}

variable "discord_token" {
  type      = string
  sensitive = true
}

variable "lambda_function_name" {
  type    = string
  default = "discord-date-decider"
}

variable "start_poll_schedule_expression" {
  type    = string
  default = "cron(0 20 15 1,2,3,4,7,8,9,10,11,12 ? *)"
}

variable "start_poll_schedule_flexible_time_window_mode" {
  type    = string
  default = "FLEXIBLE"
}

variable "start_poll_schedule_flexible_time_window_in_minutes" {
  type    = number
  default = 30
}

variable "start_poll_schedule_name" {
  type    = string
  default = "discord-date-decider-start-poll"
}

variable "end_poll_schedule_expression" {
  type    = string
  default = "cron(0 13 L 1,2,3,4,7,8,9,10,11,12 ? *)"
}

variable "end_poll_schedule_flexible_time_window_mode" {
  type    = string
  default = "FLEXIBLE"
}

variable "end_poll_schedule_flexible_time_window_in_minutes" {
  type    = number
  default = 30
}

variable "end_poll_schedule_name" {
  type    = string
  default = "discord-date-decider-end-poll"
}

variable "time_zone" {
  type    = string
  default = "Europe/Berlin"
}