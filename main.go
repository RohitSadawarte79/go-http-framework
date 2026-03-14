package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/RohitSadawarte79/go-http-framework/internal/config"
	"github.com/RohitSadawarte79/go-http-framework/internal/handler"
	"github.com/RohitSadawarte79/go-http-framework/internal/repository"
	"github.com/RohitSadawarte79/go-http-framework/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func connectDB(cfg *config.Config) (*sql.DB, error) {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping: The database is not seems to working: %w", err)
	}

	return db, nil
}

func main() {
	cfg := config.Load()

	db, err := connectDB(cfg)

	if err != nil {
		panic(err)
	}

	repo := repository.NewPostgresUserRepository(db)

	userService := service.NewUserService(repo)

	userHandler := handler.NewUserHandler(userService)

	router := NewRouter()

	router.HandleFunc("GET", "/user", userHandler.List)
	router.HandleFunc("POST", "/user", userHandler.Create)
	router.HandleFunc("GET", "/user/:id", userHandler.GetByID)

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
	err1 := http.ListenAndServe(":"+cfg.Port, stack)

	if err1 != nil {
		fmt.Println("Error starting server:", err1)
	}
}
