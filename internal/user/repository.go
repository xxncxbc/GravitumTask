package user

import (
	"GravitumTask/pkg/db"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IUserRepository interface {
	Create(newUser *User) (*User, error)
	GetById(id uint) (*User, error)
	Update(updUser *User, id uint) (*User, error)
}

type UserRepository struct {
	database *db.Db
}

func NewUserRepository(database *db.Db) *UserRepository {
	return &UserRepository{database}
}

func (repo *UserRepository) Create(newUser *User) (*User, error) {
	pass, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	newUser.Password = string(pass)
	result := repo.database.DB.Create(newUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return newUser, nil
}

func (repo *UserRepository) Update(updUser *User, id uint) (*User, error) {
	var existedUser User
	result := repo.database.DB.First(&existedUser, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	if existedUser.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	if bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(updUser.Password)) != nil {
		return nil, errors.New(ErrWrongPassword)
	}
	updUser.ID = id
	updUser.Password = existedUser.Password
	result = repo.database.DB.Clauses(clause.Returning{}).Updates(updUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return updUser, nil
}

func (repo *UserRepository) GetById(id uint) (*User, error) {
	var existedUser User
	result := repo.database.DB.First(&existedUser, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	if existedUser.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &existedUser, nil
}
