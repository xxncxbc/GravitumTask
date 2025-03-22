package user

import (
	"GravitumTask/pkg/req"
	"GravitumTask/pkg/res"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type UserHandler struct {
	Repository *UserRepository
}

func NewUserHandler(router *http.ServeMux, repository *UserRepository) {
	handler := &UserHandler{Repository: repository}
	router.HandleFunc("POST /users", handler.Create())
	router.HandleFunc("PUT /users/{id}", handler.Update())
}

func (handler *UserHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := req.HandleBody[UserCreateRequest](&w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		createdUser, err := handler.Repository.Create(&User{
			Email:    payload.Email,
			Password: payload.Password,
			Name:     payload.Name,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp := UserResponse{
			ID:        createdUser.ID,
			CreatedAt: createdUser.CreatedAt,
			UpdatedAt: createdUser.UpdatedAt,
			Name:      createdUser.Name,
			Email:     createdUser.Email,
		}
		res.Json(w, resp, http.StatusCreated)
	}
}

func (handler *UserHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		if idString == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseUint(idString, 10, 32)
		payload, err := req.HandleBody[UserUpdateRequest](&w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		updUser, err := handler.Repository.Update(&User{
			Name:     payload.Name,
			Email:    payload.Email,
			Password: payload.Password,
		}, uint(id))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp := UserResponse{
			ID:        updUser.ID,
			CreatedAt: updUser.CreatedAt,
			UpdatedAt: updUser.UpdatedAt,
			Name:      updUser.Name,
			Email:     updUser.Email,
		}
		res.Json(w, resp, http.StatusOK)
	}
}
