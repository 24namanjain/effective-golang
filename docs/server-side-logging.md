# Server-Side Logging in Golang with Zap and Zerolog

This guide covers implementing robust server-side logging in Golang applications using two popular logging libraries: **Zap** and **Zerolog**. Both libraries are designed for high-performance, structured logging that's essential for production applications.

## Table of Contents

1. [Overview](#overview)
2. [Zap Logger](#zap-logger)
3. [Zerolog Logger](#zerolog-logger)
4. [Comparison](#comparison)
5. [Best Practices](#best-practices)
6. [Advanced Configuration](#advanced-configuration)
7. [Integration Examples](#integration-examples)

## Overview

### Why Structured Logging?

Structured logging provides:
- **Machine-readable logs** for easier parsing and analysis
- **Better performance** compared to traditional string-based logging
- **Consistent log format** across your application
- **Easy integration** with log aggregation systems (ELK, Splunk, etc.)
- **Contextual information** with structured fields

### Choosing Between Zap and Zerolog

| Feature | Zap | Zerolog |
|---------|-----|---------|
| Performance | Excellent | Excellent |
| API Design | More verbose, explicit | Fluent, chainable |
| Learning Curve | Moderate | Easy |
| JSON Output | Native | Native |
| Customization | High | High |
| Zero Allocation | Yes | Yes |

## Zap Logger

### Installation

```bash
go get -u go.uber.org/zap
```

### Basic Usage

```go
package main

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func main() {
    // Create a basic logger
    logger, _ := zap.NewProduction()
    defer logger.Sync() // Flush any buffered log entries

    // Basic logging
    logger.Info("Server started",
        zap.String("port", "8080"),
        zap.String("environment", "production"),
    )

    logger.Error("Database connection failed",
        zap.String("database", "postgres"),
        zap.Error(err),
    )
}
```

### Configuration Options

#### Development Logger (Human-readable)

```go
func createDevelopmentLogger() *zap.Logger {
    config := zap.NewDevelopmentConfig()
    config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
    config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    
    logger, _ := config.Build()
    return logger
}
```

#### Production Logger (JSON)

```go
func createProductionLogger() *zap.Logger {
    config := zap.NewProductionConfig()
    config.EncoderConfig.TimeKey = "timestamp"
    config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    config.EncoderConfig.StacktraceKey = "stacktrace"
    
    logger, _ := config.Build()
    return logger
}
```

#### Custom Configuration

```go
func createCustomLogger() *zap.Logger {
    config := zap.Config{
        Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
        Development: false,
        Sampling: &zap.SamplingConfig{
            Initial:    100,
            Thereafter: 100,
        },
        Encoding:         "json",
        EncoderConfig:   zap.NewProductionEncoderConfig(),
        OutputPaths:     []string{"stdout", "/var/log/app.log"},
        ErrorOutputPaths: []string{"stderr"},
    }
    
    logger, _ := config.Build()
    return logger
}
```

### Advanced Usage

#### Structured Logging with Context

```go
type RequestLogger struct {
    logger *zap.Logger
}

func (rl *RequestLogger) LogRequest(method, path string, statusCode int, duration time.Duration) {
    rl.logger.Info("HTTP request completed",
        zap.String("method", method),
        zap.String("path", path),
        zap.Int("status_code", statusCode),
        zap.Duration("duration", duration),
        zap.String("user_agent", userAgent),
        zap.String("ip", clientIP),
    )
}

func (rl *RequestLogger) LogError(err error, context map[string]interface{}) {
    fields := make([]zap.Field, 0, len(context)+1)
    fields = append(fields, zap.Error(err))
    
    for key, value := range context {
        fields = append(fields, zap.Any(key, value))
    }
    
    rl.logger.Error("Operation failed", fields...)
}
```

#### Child Loggers with Context

```go
func handleUserRequest(logger *zap.Logger, userID string) {
    // Create a child logger with user context
    userLogger := logger.With(
        zap.String("user_id", userID),
        zap.String("request_id", generateRequestID()),
    )
    
    userLogger.Info("Processing user request")
    
    // All subsequent logs will include user context
    userLogger.Info("User data retrieved",
        zap.String("data_source", "database"),
    )
}
```

## Zerolog Logger

### Installation

```bash
go get -u github.com/rs/zerolog
```

### Basic Usage

```go
package main

import (
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "os"
    "time"
)

func main() {
    // Configure global logger
    zerolog.TimeFieldFormat = time.RFC3339
    log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
    
    // Basic logging
    log.Info().
        Str("port", "8080").
        Str("environment", "production").
        Msg("Server started")
    
    log.Error().
        Str("database", "postgres").
        Err(err).
        Msg("Database connection failed")
}
```

### Configuration Options

#### Development Mode (Human-readable)

```go
func setupDevelopmentLogger() {
    zerolog.TimeFieldFormat = time.RFC3339
    log.Logger = log.Output(zerolog.ConsoleWriter{
        Out:        os.Stdout,
        TimeFormat: time.RFC3339,
        FormatLevel: func(i interface{}) string {
            return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
        },
    })
}
```

#### Production Mode (JSON)

```go
func setupProductionLogger() {
    zerolog.TimeFieldFormat = time.RFC3339
    zerolog.SetGlobalLevel(zerolog.InfoLevel)
    
    // Optional: Write to file
    file, _ := os.OpenFile("/var/log/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    multi := zerolog.MultiLevelWriter(os.Stdout, file)
    log.Logger = zerolog.New(multi).With().Timestamp().Logger()
}
```

#### Custom Configuration

```go
func createCustomZerologLogger() zerolog.Logger {
    // Create console writer for development
    consoleWriter := zerolog.ConsoleWriter{
        Out:        os.Stdout,
        TimeFormat: time.RFC3339,
        FormatLevel: func(i interface{}) string {
            return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
        },
        FormatMessage: func(i interface{}) string {
            return fmt.Sprintf("| %s |", i)
        },
    }
    
    // Create file writer for production
    fileWriter, _ := os.OpenFile("/var/log/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    
    // Multi-writer to output to both console and file
    multi := zerolog.MultiLevelWriter(consoleWriter, fileWriter)
    
    return zerolog.New(multi).With().Timestamp().Logger()
}
```

### Advanced Usage

#### Structured Logging with Context

```go
type RequestLogger struct {
    logger zerolog.Logger
}

func (rl *RequestLogger) LogRequest(method, path string, statusCode int, duration time.Duration) {
    rl.logger.Info().
        Str("method", method).
        Str("path", path).
        Int("status_code", statusCode).
        Dur("duration", duration).
        Str("user_agent", userAgent).
        Str("ip", clientIP).
        Msg("HTTP request completed")
}

func (rl *RequestLogger) LogError(err error, context map[string]interface{}) {
    event := rl.logger.Error().Err(err)
    
    for key, value := range context {
        event = event.Interface(key, value)
    }
    
    event.Msg("Operation failed")
}
```

#### Child Loggers with Context

```go
func handleUserRequest(logger zerolog.Logger, userID string) {
    // Create a child logger with user context
    userLogger := logger.With().
        Str("user_id", userID).
        Str("request_id", generateRequestID()).
        Logger()
    
    userLogger.Info().Msg("Processing user request")
    
    // All subsequent logs will include user context
    userLogger.Info().
        Str("data_source", "database").
        Msg("User data retrieved")
}
```

#### Sampling for High-Volume Logging

```go
func createSampledLogger() zerolog.Logger {
    // Sample 10% of logs at debug level
    sampled := zerolog.Sample(&zerolog.BasicSampler{N: 10})
    
    return zerolog.New(os.Stdout).
        Sample(sampled).
        With().
        Timestamp().
        Logger()
}
```

## Comparison

### Performance Comparison

Both libraries are designed for high performance, but here's a general comparison:

- **Zap**: Slightly faster for simple logging operations
- **Zerolog**: Better performance for complex structured logging
- **Both**: Support zero-allocation logging for hot paths

### API Comparison

#### Simple Logging

```go
// Zap
logger.Info("User logged in", zap.String("user_id", "123"))

// Zerolog
log.Info().Str("user_id", "123").Msg("User logged in")
```

#### Error Logging

```go
// Zap
logger.Error("Database error", zap.Error(err), zap.String("query", sql))

// Zerolog
log.Error().Err(err).Str("query", sql).Msg("Database error")
```

#### Contextual Logging

```go
// Zap
userLogger := logger.With(zap.String("user_id", "123"))

// Zerolog
userLogger := log.With().Str("user_id", "123").Logger()
```

## Best Practices

### 1. Log Levels

Use appropriate log levels:

```go
// DEBUG: Detailed information for debugging
logger.Debug("Processing request", zap.String("request_id", reqID))

// INFO: General information about application flow
logger.Info("User authenticated", zap.String("user_id", userID))

// WARN: Something unexpected happened, but the application can continue
logger.Warn("Database connection slow", zap.Duration("duration", duration))

// ERROR: Something failed, but the application can continue
logger.Error("Failed to send email", zap.Error(err))

// FATAL: Application cannot continue
logger.Fatal("Cannot connect to database", zap.Error(err))
```

### 2. Structured Fields

Always use structured fields instead of string concatenation:

```go
// ❌ Bad
logger.Info("User " + userID + " logged in from " + ip)

// ✅ Good
logger.Info("User logged in",
    zap.String("user_id", userID),
    zap.String("ip", ip),
)
```

### 3. Error Handling

Always include error context:

```go
// ❌ Bad
logger.Error("Database operation failed")

// ✅ Good
logger.Error("Database operation failed",
    zap.Error(err),
    zap.String("operation", "user_create"),
    zap.String("user_id", userID),
)
```

### 4. Performance Considerations

```go
// Use conditional logging for expensive operations
if logger.Core().Enabled(zap.DebugLevel) {
    logger.Debug("Expensive debug info",
        zap.String("data", expensiveOperation()),
    )
}

// Use sampling for high-volume logs
logger = logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
    return zapcore.NewSamplerWithOptions(core, time.Second, 100, 10)
}))
```

### 5. Request Context

```go
func loggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // Create request-specific logger
            reqLogger := logger.With(
                zap.String("request_id", generateRequestID()),
                zap.String("method", r.Method),
                zap.String("path", r.URL.Path),
                zap.String("user_agent", r.UserAgent()),
                zap.String("ip", getClientIP(r)),
            )
            
            // Add logger to request context
            ctx := context.WithValue(r.Context(), "logger", reqLogger)
            r = r.WithContext(ctx)
            
            next.ServeHTTP(w, r)
            
            reqLogger.Info("Request completed",
                zap.Duration("duration", time.Since(start)),
            )
        })
    }
}
```

## Advanced Configuration

### Log Rotation

```go
import "gopkg.in/natefinch/lumberjack.v2"

func createRotatingLogger() *zap.Logger {
    writer := &lumberjack.Logger{
        Filename:   "/var/log/app.log",
        MaxSize:    10, // megabytes
        MaxBackups: 3,
        MaxAge:     28,   // days
        Compress:   true, // compress rotated files
    }
    
    config := zap.NewProductionConfig()
    config.OutputPaths = []string{"stdout"}
    config.ErrorOutputPaths = []string{"stderr"}
    
    logger, _ := config.Build()
    
    // Add file output
    logger = logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
        return zapcore.NewTee(core, zapcore.NewCore(
            zap.NewJSONEncoder(config.EncoderConfig),
            zapcore.AddSync(writer),
            zap.InfoLevel,
        ))
    }))
    
    return logger
}
```

### Custom Encoders

```go
func createCustomEncoder() zapcore.Encoder {
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.TimeKey = "timestamp"
    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
    encoderConfig.MessageKey = "message"
    encoderConfig.LevelKey = "level"
    
    return zapcore.NewJSONEncoder(encoderConfig)
}
```

### Sampling Configuration

```go
func createSampledLogger() *zap.Logger {
    config := zap.NewProductionConfig()
    config.Sampling = &zap.SamplingConfig{
        Initial:    100,  // Log first 100 entries
        Thereafter: 100,  // Then log every 100th entry
        Tick:       time.Second,
    }
    
    logger, _ := config.Build()
    return logger
}
```

## Integration Examples

### HTTP Server with Logging

```go
package main

import (
    "net/http"
    "time"
    
    "go.uber.org/zap"
)

func main() {
    logger, _ := zap.NewProduction()
    defer logger.Sync()
    
    // Create HTTP server with logging middleware
    server := &http.Server{
        Addr:    ":8080",
        Handler: loggingMiddleware(logger)(http.DefaultServeMux),
    }
    
    logger.Info("Starting HTTP server", zap.String("addr", server.Addr))
    
    if err := server.ListenAndServe(); err != nil {
        logger.Fatal("Server failed", zap.Error(err))
    }
}

func loggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            reqLogger := logger.With(
                zap.String("method", r.Method),
                zap.String("path", r.URL.Path),
                zap.String("remote_addr", r.RemoteAddr),
                zap.String("user_agent", r.UserAgent()),
            )
            
            reqLogger.Info("Request started")
            
            // Wrap response writer to capture status code
            wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}
            
            next.ServeHTTP(wrapped, r)
            
            reqLogger.Info("Request completed",
                zap.Int("status_code", wrapped.statusCode),
                zap.Duration("duration", time.Since(start)),
            )
        })
    }
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

### Database Operations with Logging

```go
func (db *Database) CreateUser(ctx context.Context, user *User) error {
    logger := ctx.Value("logger").(*zap.Logger)
    
    logger.Info("Creating user",
        zap.String("email", user.Email),
        zap.String("username", user.Username),
    )
    
    start := time.Now()
    
    if err := db.conn.Create(user).Error; err != nil {
        logger.Error("Failed to create user",
            zap.Error(err),
            zap.String("email", user.Email),
            zap.Duration("duration", time.Since(start)),
        )
        return err
    }
    
    logger.Info("User created successfully",
        zap.Uint("user_id", user.ID),
        zap.Duration("duration", time.Since(start)),
    )
    
    return nil
}
```

### Background Jobs with Logging

```go
func (job *EmailJob) Process(ctx context.Context) error {
    logger := ctx.Value("logger").(*zap.Logger)
    
    logger.Info("Processing email job",
        zap.String("job_id", job.ID),
        zap.String("recipient", job.Recipient),
    )
    
    start := time.Now()
    
    if err := job.sendEmail(); err != nil {
        logger.Error("Failed to send email",
            zap.Error(err),
            zap.String("job_id", job.ID),
            zap.Duration("duration", time.Since(start)),
        )
        return err
    }
    
    logger.Info("Email sent successfully",
        zap.String("job_id", job.ID),
        zap.Duration("duration", time.Since(start)),
    )
    
    return nil
}
```

## Conclusion

Both Zap and Zerolog are excellent choices for server-side logging in Golang applications. Choose based on your team's preferences and specific requirements:

- **Choose Zap** if you prefer explicit, verbose APIs and maximum performance
- **Choose Zerolog** if you prefer fluent, chainable APIs and easier learning curve

Remember to:
- Use structured logging consistently
- Include relevant context in every log entry
- Configure appropriate log levels for different environments
- Implement log rotation for production deployments
- Use sampling for high-volume applications
- Include request tracing for distributed systems

This setup will provide you with robust, performant logging that scales with your application and integrates seamlessly with modern observability stacks.
