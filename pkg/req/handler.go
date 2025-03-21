package req

import (
	"GravitumTask/internal/user"
	"GravitumTask/pkg/res"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func HandleBody(w *http.ResponseWriter, r *http.Request) (*user.User, error) {
	var payload user.User
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		res.Json(*w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(payload)
	if err != nil {
		res.Json(*w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	return &payload, nil
}
