package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const requestId contextKey = "requestId"

type Middleware func(http.Handler) http.Handler

func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		handler := next
		for i := len(middlewares) - 1; i >= 0; i-- {
			middleware := middlewares[i]
			handler = middleware(handler)
		}

		return handler
	}
}

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if panicVal := recover(); panicVal != nil {
				fmt.Printf("Panic: %v", panicVal)
				http.Error(w, "500: handler panicked.", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func Logger(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request started method:", r.Method, "path:", r.URL.Path)

		next.ServeHTTP(w, r)

		fmt.Println("Request completed")

	})

}

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uniqueId := time.Now().UnixNano()

		ctx := context.WithValue(r.Context(), requestId, uniqueId)

		r = r.WithContext(ctx)

		w.Header().Set("X-Request-ID", fmt.Sprintf("%d", uniqueId))

		next.ServeHTTP(w, r)

	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type CORSConfig struct {
	AllowedOrigins map[string]bool
	AllowedMethods []string
	AllowedHeaders []string
}

func NewCORS(cfg CORSConfig) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.AllowedOrigins[r.Header.Get("Origin")] {
				w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ", "))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ", "))
				if r.Method == "OPTIONS" {
					w.WriteHeader(http.StatusOK)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
