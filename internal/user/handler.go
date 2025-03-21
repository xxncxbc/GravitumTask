package user

import (
	"GravitumTask/pkg/req"
	"GravitumTask/pkg/res"
	"net/http"
)

type UserHandler struct {
	Repository *UserRepository
}

func NewUserHandler(router *http.ServeMux, repository *UserRepository) {
	handler := &UserHandler{Repository: repository}
	router.HandleFunc("POST /users", handler.Create())
}

func (handler *UserHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		createdUser, err := handler.Repository.Create()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, createdUser, http.StatusOK)
	}
}
