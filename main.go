package main

import (
	"fmt"
	"net/http"

	"github.com/RohitSadawarte79/go-http-framework/internal/handler"
	"github.com/RohitSadawarte79/go-http-framework/internal/repository"
	"github.com/RohitSadawarte79/go-http-framework/internal/service"
)

func main() {
	repo := repository.NewMemoryUserRepository()

	userService := service.NewUserService(repo)

	userHandler := handler.NewUserHandler(userService)

	router := NewRouter()

	router.HandleFunc("GET", "/user", userHandler.List)
	router.HandleFunc("POST", "/user", userHandler.Create)
	router.HandleFunc("GET", "/user/:id", userHandler.GetByID)

	corsConfig := CORSConfig{
		AllowedOrigins: map[string]bool{
			"http://localhost:3000": true,
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}

	corsMiddleware := NewCORS(corsConfig)

	stack := Chain(corsMiddleware, Recovery, Logger, RequestId)(router)

	fmt.Println("Listening on port 8080:", "http://localhost:8080")
	err := http.ListenAndServe(":8080", stack)

	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
