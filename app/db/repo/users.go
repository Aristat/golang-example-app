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

	row := u.db.Table("users").Select("users.id, users.email, users.encrypted_password").
		Where("users.email = ?", email).
		Limit(1).
		Row()

	err := row.Scan(&user.ID, &user.Email, &user.EncryptedPassword)

	return user, err
}

func NewAuthorsRepo(db *gorm.DB) (domain.UsersRepo, func(), error) {
	a := &UsersRepo{db: db}
	return a, func() {}, nil
}
