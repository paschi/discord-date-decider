<div align="center">

# Discord Date Decider 📅 🤖

[![Build & Test](https://img.shields.io/github/actions/workflow/status/paschi/discord-date-decider/build.yml?style=for-the-badge)](https://github.com/paschi/discord-date-decider/actions/workflows/build.yml)
[![Go Report Card](https://img.shields.io/badge/Go%20Report-A+-brightgreen?style=for-the-badge&logo=go)](https://goreportcard.com/report/github.com/paschi/discord-date-decider)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/paschi/discord-date-decider?style=for-the-badge&logo=go)](https://go.dev/)

_A serverless Discord bot that helps communities decide on dates for events by automatically creating and managing polls._
</div>

## ✨ Features

- 🤖 Serverless Discord bot running on AWS Lambda
- 📊 Creates polls automatically
- 📌 Automatically pins polls to keep them visible
- 📢 Sends announcements when new polls are created
- ⏰ Runs on a schedule using AWS EventBridge Scheduler
- 🔄 Fully automated deployment with Terraform

## 🚀 Getting Started

### Prerequisites

- Go 1.24+
- AWS Account
- Discord Bot Token
- Terraform (for deployment)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/paschi/discord-date-decider.git
   cd discord-date-decider
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the project:
   ```bash
   go build -o bootstrap ./cmd/bot
   ```

4. Create a deployment package:
   ```bash
   zip lambda-handler.zip bootstrap
   ```

## 🔧 Configuration

### Discord Bot Setup

1. Create a Discord application at the [Discord Developer Portal](https://discord.com/developers/applications)
2. Create a bot for your application
3. Enable the necessary intents for your bot
4. Copy your bot token for use in deployment

### Terraform Deployment

1. Navigate to the terraform directory:
   ```bash
   cd terraform
   ```

2. Initialize Terraform:
   ```bash
   terraform init
   ```

3. Create a `terraform.tfvars` file with your configuration:
   ```hcl
   discord_token = "your-discord-bot-token"
   discord_poll_channel = "your-poll-channel-id"
   discord_announcement_channel = "your-announcement-channel-id"
   start_poll_schedule_expression = "cron(0 0 1 * ? *)" # Run at midnight on the 1st of every month
   ```

4. Apply the Terraform configuration:
   ```bash
   terraform apply
   ```

## 🧪 Testing

Run the test suite:

```bash
go test -v ./...
```

## 📖 Usage

Once deployed, the bot will:

1. Run automatically according to the schedule you configured
2. Create a poll in the configured poll channel
3. Pin the poll to the configured poll channel
4. Send an announcement to the configured announcement channel

You can also manually trigger the function through the AWS Console or CLI.

## 🛠️ Development

### Project Structure

- `cmd/bot/` - Main application entry point
- `internal/discord/` - Discord API integration
- `internal/message/` - Message handling
- `internal/poll/` - Poll creation and management
- `terraform/` - Infrastructure as code

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
