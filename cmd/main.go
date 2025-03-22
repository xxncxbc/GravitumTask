package main

import (
	"GravitumTask/configs"
	"GravitumTask/internal/user"
	"GravitumTask/pkg/db"
	"GravitumTask/pkg/middleware"
	"net/http"
)

func App() http.Handler {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)
	router := http.NewServeMux()
	userRepository := user.NewUserRepository(database)
	user.NewUserHandler(router, userRepository)
	stack := middleware.Chain(middleware.CORS, middleware.Logging)
	return stack(router)
}

func main() {
	app := App()
	server := http.Server{
		Addr:    ":8080",
		Handler: app,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
