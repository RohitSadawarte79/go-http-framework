package main

import (
	"fmt"
	"net/http"

	"github.com/RohitSadawarte79/go-http-framework/internal/config"
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

	cfg := config.Load()
	allowedOrgins := make(map[string]bool)

	for _, origins := range cfg.AllowedOrgins {
		allowedOrgins[origins] = true
	}

	corsConfig := CORSConfig{
		AllowedOrigins: allowedOrgins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}

	corsMiddleware := NewCORS(corsConfig)

	stack := Chain(corsMiddleware, Recovery, Logger, RequestId)(router)

	fmt.Println("Listening on port ", cfg.Port, " http://localhost:", cfg.Port)
	err := http.ListenAndServe(":"+cfg.Port, stack)

	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
