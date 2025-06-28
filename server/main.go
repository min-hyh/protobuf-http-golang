package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	discoverservicepb "protobuf-http-golang/pb"
)

// customHeaderMatcher is a function that determines which HTTP headers should be forwarded as gRPC metadata
func customHeaderMatcher(key string) (string, bool) {
	// Convert HTTP header names to gRPC metadata keys
	// gRPC metadata keys are typically lowercase
	switch strings.ToLower(key) {
	case "x-custom-header-id":
		return "x-custom-header-id", true
	case "authorization":
		return "authorization", true
	case "content-type":
		return "content-type", true
	default:
		// Return false for headers we don't want to forward
		return "", false
	}
}

func main() {
	// Create a new HTTP server mux with custom options
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(customHeaderMatcher),
	)

	// Create the service implementation
	discoverService := &server{}

	// Register the HTTP handlers directly with the server implementation
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register the HTTP handlers using the direct server registration
	if err := discoverservicepb.RegisterDiscoverServiceHandlerServer(ctx, mux, discoverService); err != nil {
		log.Fatalf("Failed to register HTTP handlers: %v", err)
	}

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("Starting HTTP server on %s", httpServer.Addr)
		log.Printf("API endpoints:")
		log.Printf("  GET  /v1/get-param-in-body/{id}")
		log.Printf("  GET  /v1/get-param-in-header")
		log.Printf("  POST /v1/post/unstructured-data")
		log.Printf("  Swagger UI: http://localhost:8081/swagger-ui/")

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to serve HTTP: %v", err)
		}
	}()

	// Serve Swagger UI
	go SwaggerUI()

	// Wait for interrupt signal to gracefully shutdown the servers
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")

	// Graceful shutdown
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("Servers stopped gracefully")
}
