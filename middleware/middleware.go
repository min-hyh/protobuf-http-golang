package middleware

import (
	"log"
	"net/http"

	"google.golang.org/grpc/metadata"
)

// HeaderMiddleware is an HTTP middleware that processes headers and adds them to the context
func HeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract headers we're interested in
		log.Print("Running middleware")
		headers := make(map[string]string)

		// Process custom headers
		if customID := r.Header.Get("X-Custom-Header-Id"); customID != "" {
			log.Print("Header detected")
			headers["X-Custom-Header-Id"] = customID
		}

		// Add more header processing here as needed
		// if customHeader := r.Header.Get("X-Another-Header"); customHeader != "" {
		//     headers["X-Another-Header"] = customHeader
		// }

		// Create metadata from headers
		md := metadata.New(headers)

		// Create new context with metadata
		ctx := metadata.NewIncomingContext(r.Context(), md)

		// Create new request with updated context
		r = r.WithContext(ctx)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
