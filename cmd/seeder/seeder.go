package main

import (
	"log"

	"github.com/dinopuguh/mycap-backend/database"
	"github.com/dinopuguh/mycap-backend/seed"
)

func main() {
	if err := database.Connect(); err != nil {
		panic("Can't connect database.")
	}

	for _, seeder := range seed.AllTypes() {
		if err := seeder.Run(database.DBConn); err != nil {
			log.Fatalln(err.Error())
		}
	}
}
