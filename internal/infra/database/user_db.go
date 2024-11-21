package database

import (
	"github.com/FreitasGabriel/fullcycle-api/internal/entity"
	"gorm.io/gorm"
)

type User struct {
	DB *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{DB: db}
}

func (u *User) Create(user *entity.User) error {
	if err := u.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := u.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
