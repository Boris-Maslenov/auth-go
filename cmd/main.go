package main

import (
	"auth-test/internal/config"
	"auth-test/internal/repository/psql"
	"auth-test/internal/service"
	"auth-test/internal/transport/http/auth"
	"auth-test/internal/transport/http/middleware"
	"auth-test/internal/transport/http/server"
	"auth-test/internal/transport/http/user"
	"auth-test/pkg/hash"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const jwtSecret = "hello_go"

func main() {
	config := config.Load()
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s", config.Host, config.Port, config.Username, config.DBName, config.SSLMode, config.Password))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("DB CONNECTED")

	hasher := hash.NewSHA1Hasher("Go1")

	usersRepo := psql.NewUsers(db)
	tokensRepo := psql.NewTokens(db)
	authService := service.NewAuthService(usersRepo, tokensRepo, hasher, jwtSecret)
	authMW := middleware.Auth(authService)
	authHandler := auth.NewHandler(authService)

	mux := http.NewServeMux()
	auth.RegisterAuthRoutes(mux, authHandler)
	user.RegisterUserRoutes(mux, authMW)

	httpServer := server.NewServer(":80", mux)
	err = httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
