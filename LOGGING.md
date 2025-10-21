# Logging Implementation

This project uses **zap** (uber-go/zap) for structured, high-performance logging with multiple log levels and request body logging.

## Features

- **Multiple Log Levels**: DEBUG, INFO, WARN, ERROR
- **Request Body Logging**: Automatically logs request bodies for non-GET requests
- **Structured JSON Logging**: All logs are output in JSON format
- **Contextual Information**: Logs include timestamp, method, path, status code, latency, client IP, and user agent
- **Configurable via Environment**: Set log level via `LOG_LEVEL` environment variable

## Configuration

### Environment Variable

Set the log level in your `.env` file:

```env
LOG_LEVEL=info
```

Available levels:
- `debug` - Most verbose, includes all logs
- `info` - General information (default)
- `warn` - Warning messages
- `error` - Error messages only

### Log Output Format

Logs are output in JSON format with the following structure:

```json
{
  "level": "INFO",
  "timestamp": "2025-10-21T10:30:45.123Z",
  "caller": "handler/auth.go:41",
  "msg": "User registration attempt",
  "email": "user@example.com",
  "name": "John Doe"
}
```

## Middleware Logging

The logger middleware automatically logs all HTTP requests with:

- HTTP method
- Request path
- Query parameters
- Status code
- Request latency
- Client IP
- User agent
- Request body (for POST, PUT, PATCH, DELETE)
- Errors (if any)

### Log Level by Status Code

The middleware automatically sets the log level based on HTTP status:

- **ERROR (500+)**: Server errors
- **WARN (400-499)**: Client errors
- **INFO (200-399)**: Successful requests and redirects

### Skip Paths

By default, the following paths are not logged:
- `/health`
- `/ping`

To customize skip paths, modify `DefaultLoggerConfig` in `internal/middleware/logger.go`.

## Usage in Handlers

### Basic Logging

```go
h.logger.Info("User created successfully",
    zap.Uint("user_id", user.ID),
    zap.String("email", user.Email),
)
```

### Log Levels

```go
// Debug - detailed information for debugging
h.logger.Debug("Fetching user profile",
    zap.Any("user_id", userID),
)

// Info - general informational messages
h.logger.Info("Login successful",
    zap.Uint("user_id", user.ID),
)

// Warn - warning messages (recoverable issues)
h.logger.Warn("Login failed - invalid password",
    zap.String("email", req.Email),
)

// Error - error messages (serious issues)
h.logger.Error("Failed to create user",
    zap.String("error", err.Error()),
    zap.String("email", req.Email),
)
```

### Structured Fields

Use typed fields for better performance and type safety:

```go
zap.String("key", stringValue)
zap.Int("key", intValue)
zap.Uint("key", uintValue)
zap.Int64("key", int64Value)
zap.Bool("key", boolValue)
zap.Duration("key", durationValue)
zap.Time("key", timeValue)
zap.Any("key", anyValue)  // Use sparingly
```

## Example Request Log

### Request

```bash
POST /api/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "secret123"
}
```

### Log Output

```json
{
  "level": "INFO",
  "timestamp": "2025-10-21T10:30:45.123Z",
  "caller": "middleware/logger.go:99",
  "msg": "Request completed",
  "method": "POST",
  "path": "/api/auth/register",
  "query": "",
  "status": 201,
  "latency": "15.234ms",
  "client_ip": "127.0.0.1",
  "user_agent": "PostmanRuntime/7.32.3",
  "request_body": "{\"name\":\"John Doe\",\"email\":\"john@example.com\",\"password\":\"secret123\"}"
}
```

## Best Practices

1. **Use Appropriate Log Levels**
   - DEBUG: Detailed diagnostic information
   - INFO: Normal operations and business events
   - WARN: Unexpected but recoverable issues
   - ERROR: Serious problems that need attention

2. **Add Context**
   - Include relevant IDs (user_id, order_id, etc.)
   - Add error messages from caught errors
   - Include request identifiers for tracing

3. **Avoid Sensitive Data**
   - Never log passwords or tokens
   - Be careful with personal information
   - Consider masking sensitive fields

4. **Use Structured Fields**
   - Prefer typed fields over string formatting
   - Makes logs easier to search and analyze
   - Better performance

5. **Production Settings**
   - Set `LOG_LEVEL=info` or `LOG_LEVEL=warn` in production
   - Use `LOG_LEVEL=debug` only for troubleshooting
   - Consider log aggregation tools (ELK, Datadog, etc.)

## Testing Logs

You can test different log levels by setting the environment variable:

```bash
# Show all logs including debug
LOG_LEVEL=debug go run cmd/server/main.go

# Show only info, warn, and error logs (default)
LOG_LEVEL=info go run cmd/server/main.go

# Show only warnings and errors
LOG_LEVEL=warn go run cmd/server/main.go

# Show only errors
LOG_LEVEL=error go run cmd/server/main.go
```
