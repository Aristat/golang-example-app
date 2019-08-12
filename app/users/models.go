package users

import (
	"errors"

	"github.com/jinzhu/gorm"

	"github.com/aristat/golang-gin-oauth2-example-app/common"
)

var (
	passwordIsEmpty = errors.New("10001 email is not valid")
)

type UserModel struct {
	ID                int
	Email             string
	EncryptedPassword string
}

func (u *UserModel) setPassword(password string, cost int) error {
	if len(password) == 0 {
		return passwordIsEmpty
	}

	passwordHash, _ := common.HashPassword(password, cost)
	u.EncryptedPassword = string(passwordHash)
	return nil
}

func FindByEmail(db *gorm.DB, email string) (*UserModel, error) {
	u := &UserModel{}

	row := db.Table("users").Select("users.id, users.email, users.encrypted_password").
		Where("users.email = ?", email).
		Limit(1).
		Row()

	err := row.Scan(&u.ID, &u.Email, &u.EncryptedPassword)

	return u, err
}
