package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error   string            `json:"error"`
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// ErrorHandler defines the interface for custom error handling logic
type ErrorHandler interface {
	HandleError(ctx context.Context, err error, req *http.Request) *ErrorResponse
}

// DefaultErrorHandler provides default error handling logic
type DefaultErrorHandler struct{}

// HandleError implements the default error handling logic
func (h *DefaultErrorHandler) HandleError(ctx context.Context, err error, req *http.Request) *ErrorResponse {
	// Convert gRPC status to HTTP status
	grpcStatus, ok := status.FromError(err)
	if ok {
		return h.handleGRPCStatus(grpcStatus, req)
	}

	// Handle runtime errors
	if runtimeErr, ok := err.(*runtime.HTTPStatusError); ok {
		return h.handleRuntimeError(runtimeErr, req)
	}

	// Handle panic errors
	if panicErr, ok := err.(*PanicError); ok {
		return h.handlePanicError(panicErr, req)
	}

	// Default error response
	return &ErrorResponse{
		Error:   "Internal Server Error",
		Code:    http.StatusInternalServerError,
		Message: "An unexpected error occurred",
		Details: map[string]string{
			"request_path": req.URL.Path,
			"method":       req.Method,
		},
	}
}

// handleGRPCStatus handles gRPC status errors
func (h *DefaultErrorHandler) handleGRPCStatus(grpcStatus *status.Status, req *http.Request) *ErrorResponse {
	httpStatus := grpcStatusToHTTPStatus(grpcStatus.Code())

	response := &ErrorResponse{
		Error:   grpcStatus.Code().String(),
		Code:    httpStatus,
		Message: grpcStatus.Message(),
		Details: map[string]string{
			"request_path": req.URL.Path,
			"method":       req.Method,
		},
	}

	// Add additional details based on error code
	switch grpcStatus.Code() {
	case codes.InvalidArgument:
		response.Message = "Invalid request parameters"
	case codes.NotFound:
		response.Message = "Resource not found"
	case codes.PermissionDenied:
		response.Message = "Access denied"
	case codes.Unauthenticated:
		response.Message = "Authentication required"
	}

	return response
}

// handleRuntimeError handles runtime.HTTPStatusError
func (h *DefaultErrorHandler) handleRuntimeError(err *runtime.HTTPStatusError, req *http.Request) *ErrorResponse {
	return &ErrorResponse{
		Error:   "HTTP Status Error",
		Code:    err.HTTPStatus,
		Message: err.Err.Error(),
		Details: map[string]string{
			"request_path": req.URL.Path,
			"method":       req.Method,
		},
	}
}

// handlePanicError handles panic errors
func (h *DefaultErrorHandler) handlePanicError(err *PanicError, req *http.Request) *ErrorResponse {
	return &ErrorResponse{
		Error:   "Panic Error",
		Code:    http.StatusInternalServerError,
		Message: "A panic occurred while processing the request",
		Details: map[string]string{
			"request_path": req.URL.Path,
			"method":       req.Method,
			"panic_msg":    err.Error(),
			"stack_trace":  err.StackTrace,
		},
	}
}

// PanicError represents a panic error with stack trace
type PanicError struct {
	Message    string
	StackTrace string
}

func (e *PanicError) Error() string {
	return e.Message
}

// grpcStatusToHTTPStatus converts gRPC status codes to HTTP status codes
func grpcStatusToHTTPStatus(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// ErrorHandlingMiddleware creates a middleware that handles errors with custom logic
func ErrorHandlingMiddleware(errorHandler ErrorHandler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a custom response writer that captures errors
			errorWriter := &ErrorResponseWriter{
				ResponseWriter: w,
				ErrorHandler:   errorHandler,
				Request:        r,
			}

			// Recover from panics
			defer func() {
				if rec := recover(); rec != nil {
					stackTrace := string(debug.Stack())
					panicErr := &PanicError{
						Message:    fmt.Sprintf("Panic: %v", rec),
						StackTrace: stackTrace,
					}

					log.Printf("Panic recovered: %v\n%s", rec, stackTrace)

					response := errorHandler.HandleError(r.Context(), panicErr, r)
					writeErrorResponse(w, response)
				}
			}()

			next.ServeHTTP(errorWriter, r)
		})
	}
}

// ErrorResponseWriter wraps http.ResponseWriter to capture errors
type ErrorResponseWriter struct {
	http.ResponseWriter
	ErrorHandler ErrorHandler
	Request      *http.Request
	statusCode   int
	written      bool
}

// WriteHeader captures the status code
func (w *ErrorResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
	w.written = true
}

// Write captures the response
func (w *ErrorResponseWriter) Write(data []byte) (int, error) {
	w.written = true
	return w.ResponseWriter.Write(data)
}

// writeErrorResponse writes an error response in JSON format
func writeErrorResponse(w http.ResponseWriter, response *ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Code)

	jsonData, err := json.Marshal(response)
	if err != nil {
		// Fallback to simple error response
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal Server Error","code":500,"message":"Failed to serialize error response"}`))
		return
	}

	w.Write(jsonData)
}

// CustomErrorHandler is an example of a custom error handler
type CustomErrorHandler struct {
	DefaultErrorHandler
	// Add custom fields as needed
	LogErrors bool
}

// HandleError implements custom error handling logic
func (h *CustomErrorHandler) HandleError(ctx context.Context, err error, req *http.Request) *ErrorResponse {
	// Log errors if enabled
	if h.LogErrors {
		log.Printf("Error handling request %s %s: %v", req.Method, req.URL.Path, err)
	}

	// Call the default handler first
	response := h.DefaultErrorHandler.HandleError(ctx, err, req)

	// Add custom logic here
	// For example, you could:
	// - Add request ID to error details
	// - Send errors to external monitoring services
	// - Add custom error codes
	// - Modify error messages based on environment

	// Example: Add request ID if available
	if requestID := req.Header.Get("X-Request-ID"); requestID != "" {
		if response.Details == nil {
			response.Details = make(map[string]string)
		}
		response.Details["request_id"] = requestID
	}

	return response
}
