package user

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
