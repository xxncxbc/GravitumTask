package main

import (
	"GravitumTask/configs"
	"GravitumTask/internal/user"
	"GravitumTask/pkg/db"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)
	router := http.NewServeMux()
	userRepository := user.NewUserRepository(database)
	user.NewUserHandler(router, userRepository)
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
