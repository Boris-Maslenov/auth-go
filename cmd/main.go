package main

import (
	"auth-test/internal/repository/psql"
	"auth-test/internal/service"
	"auth-test/internal/transport/http/auth"
	"auth-test/internal/transport/http/server"
	"auth-test/pkg/hash"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

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
	authService := service.NewAuthService(hasher, usersRepo)
	authHandler := auth.NewHandler(authService)

	router := server.NewRouter(authHandler)
	httpServer := server.NewServer(":80", router)
	err = httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
