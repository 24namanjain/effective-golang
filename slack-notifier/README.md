# Slack Notifier

A modular Go application for sending notifications to Slack with support for multiple event types, severity levels, and rich formatting.

## 🚀 Features

- **Modular Architecture**: Clean separation of concerns with configurable components
- **Multiple Event Types**: Support for system, user, business, and alert events
- **Severity Levels**: Info, Warning, Error, and Critical notifications
- **Rich Slack Formatting**: Beautiful message blocks with metadata and context
- **Concurrent Processing**: Multi-worker architecture for high throughput
- **Graceful Shutdown**: Proper cleanup and final notifications
- **Configuration Management**: Environment-based configuration with validation
- **Error Handling**: Comprehensive error handling and retry logic

## 📋 Prerequisites

- Go 1.21 or higher
- Slack Bot Token (required)
- Slack App configured with appropriate permissions

## 🏗️ Project Structure

```
slack-notifier/
├── cmd/
│   └── server/
│       └── main.go              # Main application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── events/
│   │   └── event.go             # Event types and builders
│   └── notifier/
│       └── notifier.go          # Core notification service
├── pkg/
│   └── slack/
│       └── client.go            # Slack API client wrapper
├── examples/
│   ├── basic_usage.go           # Basic usage example
│   └── event_types.go           # All event types demonstration
├── docs/                        # Documentation
├── go.mod                       # Go module file
└── README.md                    # This file
```

## 🔧 Configuration

### Environment Variables

The application uses the following environment variables:

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `SLACK_BOT_TOKEN` | Slack Bot User OAuth Token | - | ✅ Yes |
| `SLACK_CLIENT_ID` | Slack App Client ID | `9371894822341.9375537877046` | No |
| `SLACK_APP_ID` | Slack App ID | `A09B1FTRT1C` | No |
| `SLACK_CLIENT_SECRET` | Slack App Client Secret | `1b5dc7da4540da29a5e02ec9bbaf69e5` | No |
| `SLACK_CHANNEL` | Default Slack channel | `#general` | No |
| `ENVIRONMENT` | Application environment | `development` | No |

### Setting up Slack Bot Token

1. Go to [Slack API Apps](https://api.slack.com/apps)
2. Create a new app or select an existing one
3. Go to "OAuth & Permissions"
4. Add the following bot token scopes:
   - `chat:write` - Send messages to channels
   - `chat:write.public` - Send messages to public channels
   - `channels:read` - Read channel information
5. Install the app to your workspace
6. Copy the "Bot User OAuth Token" (starts with `xoxb-`)

## 🚀 Quick Start

### 1. Set up environment variables

```bash
export SLACK_BOT_TOKEN="xoxb-your-bot-token-here"
export SLACK_CHANNEL="#your-channel"
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Run the main application

```bash
go run cmd/server/main.go
```

### 4. Run examples

```bash
# Basic usage example
go run examples/basic_usage.go

# All event types demonstration
go run examples/event_types.go
```

## 📚 Usage Examples

### Basic Message

```go
service, err := notifier.NewNotifierService(config, 2)
if err != nil {
    log.Fatal(err)
}

// Send a simple text message
err = service.SendMessage("Hello from Slack Notifier! 👋")
```

### Event with Metadata

```go
event := events.NewEvent(events.EventTypeUserLogin).
    WithTitle("User Login").
    WithMessage("User john@example.com has logged in").
    WithSeverity(events.SeverityInfo).
    WithUserID("john@example.com").
    WithMetadata("ip_address", "192.168.1.100").
    WithMetadata("user_agent", "Chrome/91.0").
    Build()

service.SendEvent(event)
```

### Critical Alert

```go
alertEvent := events.NewEvent(events.EventTypeServiceDown).
    WithTitle("Service Down Alert").
    WithMessage("Database service is not responding").
    WithSeverity(events.SeverityCritical).
    WithMetadata("service", "database").
    WithMetadata("downtime", "5 minutes").
    Build()

service.SendEvent(alertEvent)
```

## 🎯 Event Types

### System Events
- `system_startup` - Application startup
- `system_shutdown` - Application shutdown
- `system_error` - System errors

### User Events
- `user_login` - User login
- `user_logout` - User logout
- `user_registration` - New user registration

### Business Events
- `order_created` - New order created
- `order_completed` - Order completed
- `order_cancelled` - Order cancelled
- `payment_received` - Payment received
- `payment_failed` - Payment failed

### Alert Events
- `high_cpu_usage` - High CPU usage alert
- `high_memory_usage` - High memory usage alert
- `disk_space_low` - Low disk space alert
- `service_down` - Service down alert

## 🎨 Severity Levels

- **Info** (ℹ️) - General information
- **Warning** (⚠️) - Warning conditions
- **Error** (❌) - Error conditions
- **Critical** (🚨) - Critical issues requiring immediate attention

## 🔧 Architecture

### Components

1. **Config Module** (`internal/config/`)
   - Loads and validates configuration
   - Supports environment variables and .env files
   - Provides helper methods for environment detection

2. **Events Module** (`internal/events/`)
   - Defines event types and structures
   - Provides fluent builder pattern for event creation
   - Handles event ID generation and metadata

3. **Slack Client** (`pkg/slack/`)
   - Wraps Slack API client
   - Handles message formatting and blocks
   - Provides connection testing and channel info

4. **Notifier Service** (`internal/notifier/`)
   - Orchestrates event processing
   - Manages worker goroutines
   - Handles graceful shutdown
   - Provides statistics and monitoring

### Concurrency Model

- **Event Queue**: Buffered channel for event processing
- **Worker Pool**: Configurable number of worker goroutines
- **Graceful Shutdown**: Context cancellation and wait groups
- **Error Handling**: Separate error notification channel

## 🧪 Testing

### Run tests

```bash
go test ./...
```

### Test specific components

```bash
# Test configuration
go test ./internal/config

# Test events
go test ./internal/events

# Test Slack client
go test ./pkg/slack
```

## 📊 Monitoring

The notifier service provides statistics:

```go
stats := service.GetStats()
// Returns: map[string]interface{}{
//   "queue_size": 5,
//   "workers": 3,
//   "active": true,
// }
```

## 🔒 Security Considerations

1. **Token Security**: Never commit Slack tokens to version control
2. **Environment Variables**: Use environment variables for sensitive data
3. **Channel Permissions**: Ensure the bot has appropriate channel permissions
4. **Rate Limiting**: The service respects Slack API rate limits

## 🚨 Troubleshooting

### Common Issues

1. **"SLACK_BOT_TOKEN is required"**
   - Set the `SLACK_BOT_TOKEN` environment variable
   - Ensure the token starts with `xoxb-`

2. **"Failed to authenticate with Slack"**
   - Verify the bot token is correct
   - Check if the bot is installed in your workspace
   - Ensure the bot has the required scopes

3. **"Failed to send event to Slack"**
   - Check if the bot is in the target channel
   - Verify channel permissions
   - Check Slack API status

### Debug Mode

Enable debug logging by setting the log level:

```bash
export LOG_LEVEL=debug
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- [slack-go/slack](https://github.com/slack-go/slack) - Official Slack Go SDK
- [joho/godotenv](https://github.com/joho/godotenv) - Environment variable loading
- [stretchr/testify](https://github.com/stretchr/testify) - Testing utilities
