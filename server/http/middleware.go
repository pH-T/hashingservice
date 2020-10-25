package http

import (
	"context"
	"log"
	"net/http"
	"time"
)

// GETMiddleware enforces the usage of HTTP GET
func GETMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			reqid := r.Context().Value("reqid")
			log.Printf("[%s] Error: %s", reqid, "expected GET")

			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// POSTMiddleware enforces the usage of HTTP POST
func POSTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			reqid := r.Context().Value("reqid")
			log.Printf("[%s] Error: %s", reqid, "expected POST")

			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware adds basic request logging and adds a requestID to the context of the request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqid := randString(15)
		ctx := r.Context()
		ctx = context.WithValue(ctx, "reqid", reqid)

		log.Printf("[%s] %s %s %s", reqid, r.RemoteAddr, r.RequestURI, r.UserAgent())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// TimeMiddleware measures the time it took to process the request
func TimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqid := r.Context().Value("reqid")
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf("[%s] %v", reqid, time.Now().Sub(start))
	})
}
