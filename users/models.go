package users

import (
	"github.com/aristat/golang-gin-oauth2-example-app/common"
)

type UserModel struct {
	ID                int
	Email             string
	EncryptedPassword string
}

func (u *UserModel) setPassword(password string) error {
	if len(password) == 0 {
		return common.PasswordIsEmpty
	}

	passwordHash, _ := common.HashPassword(password)
	u.EncryptedPassword = string(passwordHash)
	return nil
}

func FindByEmail(env *common.Env, email string) (*UserModel, error) {
	db := env.DB
	u := &UserModel{}

	err := db.
		QueryRow(`SELECT id, email, encrypted_password FROM users WHERE email = $1 LIMIT 1`, email).
		Scan(&u.ID, &u.Email, &u.EncryptedPassword)

	return u, err
}
