package domain

type Product struct {
	ID   int    `gorm:"column:id"`
	Name string `gorm:"column:email"`
}
