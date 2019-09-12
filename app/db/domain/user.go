package domain

type User struct {
	ID                int    `gorm:"column:id"`
	Email             string `gorm:"column:email"`
	EncryptedPassword string `gorm:"column:encrypted_password"`
}

// UsersRepo interface
type UsersRepo interface {
	CreateUser(email string, password string) (*User, error)
	FindByEmail(email string) (*User, error)
}
