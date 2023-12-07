package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Username string
	Email    string
	SocialID string
}

type Order struct {
	gorm.Model
	Name     string
	CoverUrl string
	DoneAt   sql.NullTime
	AuthorID string
	Author   User
}

type Item struct {
	gorm.Model
	Quantity  float64
	Recipient string
	ProductID string
	Product   Product
}

type Product struct {
	gorm.Model
	Name    string
	Price   sql.NullFloat64
	OrderID string
	Order   Order
}
