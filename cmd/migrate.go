package main

import (
	"github.com/dinopuguh/mycap-backend/database"
	"github.com/dinopuguh/mycap-backend/migrations"
)

func main() {
	err := database.Connect()
	if err != nil {
		panic("Can't connect database.")
	}

	migrations.All()
}
