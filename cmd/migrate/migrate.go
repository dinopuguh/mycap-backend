package main

import (
	"github.com/dinopuguh/mycap-backend/database"
	"github.com/dinopuguh/mycap-backend/migrations"
)

func main() {
	if err := database.Connect(); err != nil {
		panic("Can't connect database.")
	}

	migrations.All()
}
