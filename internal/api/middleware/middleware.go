package middleware

import (
	"encoding/json"
	"loan/internal/domain"
	"log"
	"net/http"
	"time"
)

// Logger is a middleware that logs incoming requests
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		// Log the request
		log.Printf(
			"%s %s %s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

// ContentTypeJSON sets the Content-Type header to application/json
func ContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// ErrorHandler is a middleware that handles errors in a standard way
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v", err)

				resp := domain.NewErrorResponse(http.StatusInternalServerError, "Internal server error")

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(resp)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
