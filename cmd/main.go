package main

import (
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
	// todo: брать данные из сonfig env
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres dbname=users_db sslmode=disable password=73007300")
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
