package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

var jsTemplate *template.Template

func init() {
	// Parse the JavaScript template file once at startup
	var err error
	jsTemplate, err = template.ParseFiles("templates/script.js")
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}
}

func serveJS(w http.ResponseWriter, r *http.Request) {
	// Get the `due_date` from URL path parameters
	dueDateParam := chi.URLParam(r, "dueDate")
	if dueDateParam == "" {
		http.Error(w, "due_date parameter is required", http.StatusBadRequest)
		return
	}

	// Validate the date format
	_, err := time.Parse("2006-01-02", dueDateParam)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD.", http.StatusBadRequest)
		return
	}

	// Set the Content-Type to application/javascript
	w.Header().Set("Content-Type", "application/javascript")

	// Execute the template with the dynamic due_date
	if err := jsTemplate.Execute(w, map[string]string{"DueDate": dueDateParam}); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func main() {
	r := chi.NewRouter()

	// Add some basic middleware
	r.Use(middleware.Logger)         // Logs HTTP requests
	r.Use(middleware.Recoverer)      // Recovers from panics
	r.Use(middleware.Heartbeat("/")) // Responds with 200 OK on `/`

	// Add rate-limiting middleware (e.g., max 5 requests per second per IP)
	r.Use(httprate.LimitByIP(5, 1*time.Second))

	// Define the route for the JavaScript file with a path parameter
	r.Get("/{dueDate}", serveJS)

	// Start the server
	log.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
