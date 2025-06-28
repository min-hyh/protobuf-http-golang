package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	discoverservicepb "protobuf-http-golang/pb"
)

func main() {
	// Create a new HTTP server mux
	mux := runtime.NewServeMux()

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
