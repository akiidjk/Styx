package dbstore

import (
	"github.com/akiidjk/styx/internal/web/internal/hash"
	"github.com/akiidjk/styx/internal/web/internal/store"

	"gorm.io/gorm"
)

type UserStore struct {
	db           *gorm.DB
	passwordhash hash.PasswordHash
}

type NewUserStoreParams struct {
	DB           *gorm.DB
	PasswordHash hash.PasswordHash
}

func NewUserStore(params NewUserStoreParams) *UserStore {
	return &UserStore{
		db:           params.DB,
		passwordhash: params.PasswordHash,
	}
}

func (s *UserStore) CreateUser(email string, password string) error {

	hashedPassword, err := s.passwordhash.GenerateFromPassword(password)
	if err != nil {
		return err
	}

	return s.db.Create(&store.User{
		Email:    email,
		Password: hashedPassword,
	}).Error
}

func (s *UserStore) GetUser(email string) (*store.User, error) {

	var user store.User
	err := s.db.Where("email = ?", email).First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, err
}
