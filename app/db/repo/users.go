package repo

import (
	"github.com/aristat/golang-example-app/app/db/domain"
	"github.com/jinzhu/gorm"
)

type UsersRepo struct {
	db *gorm.DB
}

func (u *UsersRepo) CreateUser(email string, encryptPassword string) (*domain.User, error) {
	user := &domain.User{Email: email, EncryptedPassword: encryptPassword}

	if err := u.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UsersRepo) FindByEmail(email string) (*domain.User, error) {
	user := &domain.User{}

	err := u.db.
		Table("users").
		Select("id, email, encrypted_password").
		Where("users.email = ?", email).
		Limit(1).
		Scan(&user).Error

	return user, err
}

func NewUsersRepo(db *gorm.DB) (domain.UsersRepo, func(), error) {
	a := &UsersRepo{db: db}
	return a, func() {}, nil
}
