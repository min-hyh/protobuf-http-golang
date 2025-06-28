package main

import (
	"log"
	"net/http"
)

func SwaggerUI() {
	swaggerMux := http.NewServeMux()
	swaggerMux.HandleFunc("/swagger-ui/", func(w http.ResponseWriter, r *http.Request) {
		// Serve the swagger.json file
		if r.URL.Path == "/swagger-ui/swagger.json" {
			http.ServeFile(w, r, "pb/discover.swagger.json")
			return
		}

		// Serve a simple Swagger UI
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
<title>Discover Service API</title>
<link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui.css" />
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-bundle.js"></script>
<script>
	window.onload = function() {
		SwaggerUIBundle({
			url: '/swagger-ui/swagger.json',
			dom_id: '#swagger-ui',
			presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
			layout: "BaseLayout"
		});
	};
</script>
</body>
</html>
		`))
	})

	swaggerServer := &http.Server{
		Addr:    ":8081",
		Handler: swaggerMux,
	}

	log.Printf("Starting Swagger UI server on %s", swaggerServer.Addr)
	if err := swaggerServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to serve Swagger UI: %v", err)
	}
}
