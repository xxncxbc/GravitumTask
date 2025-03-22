package req

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
)

// обработка тела запроса
func HandleBody[T any](w *http.ResponseWriter, r *http.Request) (*T, error) {
	var payload T
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}
