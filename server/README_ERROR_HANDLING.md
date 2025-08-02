# Error Handling Middleware

This middleware provides comprehensive error handling for gRPC-HTTP gateway servers with custom logic support.

## Features

- **Standardized Error Responses**: All errors are returned in a consistent JSON format
- **gRPC Status Code Mapping**: Automatically converts gRPC status codes to appropriate HTTP status codes
- **Panic Recovery**: Catches and handles panics gracefully
- **Custom Error Logic**: Supports custom error handling logic through the `ErrorHandler` interface
- **Request Context**: Includes request details in error responses
- **Stack Trace Support**: Captures stack traces for panic errors

## Error Response Format

All errors are returned in the following JSON format:

```json
{
  "error": "INVALID_ARGUMENT",
  "code": 400,
  "message": "Invalid request parameters",
  "details": {
    "request_path": "/v1/get-param-in-body/test",
    "method": "GET",
    "request_id": "test-request-123"
  }
}
```

## Usage

### Basic Setup

```go
// Create a custom error handler
errorHandler := &CustomErrorHandler{
    LogErrors: true, // Enable error logging
}

// Create HTTP server with error handling middleware
httpServer := &http.Server{
    Addr: ":8080",
    Handler: ErrorHandlingMiddleware(errorHandler)(mux),
}
```

### Custom Error Handler

You can create your own error handler by implementing the `ErrorHandler` interface:

```go
type MyCustomErrorHandler struct {
    DefaultErrorHandler
    // Add your custom fields
    Environment string
    LogToExternalService bool
}

func (h *MyCustomErrorHandler) HandleError(ctx context.Context, err error, req *http.Request) *ErrorResponse {
    // Call the default handler first
    response := h.DefaultErrorHandler.HandleError(ctx, err, req)
    
    // Add your custom logic here
    if h.LogToExternalService {
        // Send error to external monitoring service
        h.sendToMonitoringService(err, req)
    }
    
    // Add environment-specific details
    if response.Details == nil {
        response.Details = make(map[string]string)
    }
    response.Details["environment"] = h.Environment
    
    return response
}
```

## Error Types Handled

### 1. gRPC Status Errors
- `codes.InvalidArgument` → HTTP 400 Bad Request
- `codes.NotFound` → HTTP 404 Not Found
- `codes.PermissionDenied` → HTTP 403 Forbidden
- `codes.Unauthenticated` → HTTP 401 Unauthorized
- `codes.Internal` → HTTP 500 Internal Server Error
- And many more...

### 2. Runtime Errors
- HTTP status errors from the gRPC gateway runtime

### 3. Panic Errors
- Catches panics and returns structured error responses
- Includes stack trace in development environments

## Example Error Scenarios

### Missing Required Fields
```bash
curl -X GET "http://localhost:8080/v1/get-param-in-body/?content=test"
```
Response:
```json
{
  "error": "INVALID_ARGUMENT",
  "code": 400,
  "message": "Invalid request parameters",
  "details": {
    "request_path": "/v1/get-param-in-body/",
    "method": "GET"
  }
}
```

### Resource Not Found
```bash
curl -X GET "http://localhost:8080/v1/get-param-in-body/not-found?content=test"
```
Response:
```json
{
  "error": "NOT_FOUND",
  "code": 404,
  "message": "Resource not found",
  "details": {
    "request_path": "/v1/get-param-in-body/not-found",
    "method": "GET"
  }
}
```

### Permission Denied
```bash
curl -X GET "http://localhost:8080/v1/get-param-in-body/unauthorized?content=test"
```
Response:
```json
{
  "error": "PERMISSION_DENIED",
  "code": 403,
  "message": "Access denied",
  "details": {
    "request_path": "/v1/get-param-in-body/unauthorized",
    "method": "GET"
  }
}
```

## Testing

Run the test script to see all error scenarios in action:

```bash
./test_errors.sh
```

This will test various error conditions and show the corresponding error responses.

## Customization

### Adding Custom Error Codes

You can extend the error handling by adding custom error codes:

```go
// In your server methods
if someCondition {
    return nil, status.Errorf(codes.InvalidArgument, "custom error message")
}
```

### Adding Request-Specific Details

The middleware automatically includes request details, but you can add more:

```go
// In your custom error handler
if userID := req.Header.Get("X-User-ID"); userID != "" {
    response.Details["user_id"] = userID
}
```

### Environment-Specific Error Messages

You can modify error messages based on the environment:

```go
func (h *CustomErrorHandler) HandleError(ctx context.Context, err error, req *http.Request) *ErrorResponse {
    response := h.DefaultErrorHandler.HandleError(ctx, err, req)
    
    // Modify messages for production
    if os.Getenv("ENVIRONMENT") == "production" {
        response.Message = "An error occurred. Please try again later."
    }
    
    return response
}
```

## Best Practices

1. **Always return gRPC status errors** from your service methods
2. **Use appropriate gRPC status codes** for different error types
3. **Include meaningful error messages** that help with debugging
4. **Add request context** to error responses for better debugging
5. **Log errors appropriately** for monitoring and debugging
6. **Handle panics gracefully** to prevent server crashes
7. **Use custom error handlers** for environment-specific logic

## Integration with Monitoring

You can integrate with external monitoring services in your custom error handler:

```go
func (h *CustomErrorHandler) HandleError(ctx context.Context, err error, req *http.Request) *ErrorResponse {
    response := h.DefaultErrorHandler.HandleError(ctx, err, req)
    
    // Send to monitoring service
    go h.sendToMonitoringService(err, req, response)
    
    return response
}
``` 