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
	Repository IUserRepository
}

func NewUserHandler(router *http.ServeMux, repository IUserRepository) {
	handler := &UserHandler{Repository: repository}
	router.HandleFunc("POST /users", handler.Create())
	router.HandleFunc("PUT /users/{id}", handler.Update())
	router.HandleFunc("GET /users/{id}", handler.Get())
}

func getIdFromPath(r *http.Request) (uint, error) {
	idString := r.PathValue("id")
	if idString == "" {
		return 0, errors.New(ErrEmptyID)
	}
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func ModelToUserResponse(user *User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Email:     user.Email,
	}
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
		resp := ModelToUserResponse(createdUser)
		res.Json(w, resp, http.StatusCreated)
	}
}

func (handler *UserHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getIdFromPath(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
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
		resp := ModelToUserResponse(updUser)
		res.Json(w, resp, http.StatusOK)
	}
}

func (handler *UserHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getIdFromPath(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := handler.Repository.GetById(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp := ModelToUserResponse(user)
		res.Json(w, resp, http.StatusOK)
	}
}
