package repo

import (
	"github.com/aristat/golang-oauth2-example-app/app/db/domain"
	"github.com/jinzhu/gorm"
)

type UsersRepo struct {
	db *gorm.DB
}

func (u *UsersRepo) FindByEmail(email string) (*domain.User, error) {
	user := &domain.User{}

	err := u.db.Table("users").Select("id, email, encrypted_password").
		Where("users.email = ?", email).
		Limit(1).
		Scan(&user).Error

	return user, err
}

func NewUsersRepo(db *gorm.DB) (domain.UsersRepo, func(), error) {
	a := &UsersRepo{db: db}
	return a, func() {}, nil
}
