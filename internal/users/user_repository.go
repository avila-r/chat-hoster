package users

import (
	"gorm.io/gorm"
)

type Repository struct {
	database *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) CreateUser(user *User) (*User, error) {
	result := r.database.Create(user)

	return user, result.Error
}

func (r *Repository) FindUserByEmail(email string) (*User, error) {
	var (
		user User
	)

	result := r.database.Where("email = ?", email).First(&user)

	err := result.Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
