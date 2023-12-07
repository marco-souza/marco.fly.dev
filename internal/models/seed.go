package models

import (
	"database/sql"
	"fmt"
	"math/rand"
)

func Seed() {
	user := User{
		Name:     "Marco",
		Username: "marco-souza",
	}

	db := Connect()
	db.Create(&user)

	name := "Compras de natal"
	order := Order{
		Name: name,
		CoverUrl: fmt.Sprintf(
			"https://source.unsplash.com/random/?%s&%d&w=%d",
			name,
			rand.Intn(100),
			250,
		),
		Author: user,
	}
	db.Create(&order)

	products := []string{
		"Presentes",
		"Comidas",
		"Refri",
		"Doces",
		"Cha",
	}
	product := Product{
		Name: products[rand.Intn(len(products))],
		Price: sql.NullFloat64{
			Float64: rand.Float64() * 100,
		},
		Order: order,
	}
	db.Create(&product)

	for range [5]int{} {
		item := Item{
			Product:   product,
			Quantity:  rand.Float64() * 10,
			Recipient: "someone else",
		}
		db.Create(&item)
	}
}
