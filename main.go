package main

import (
	"flag"
	"log"
	"os"

	"github.com/dinopuguh/mycap-backend/database"
	_ "github.com/dinopuguh/mycap-backend/docs"
	"github.com/dinopuguh/mycap-backend/migrations"
	"github.com/dinopuguh/mycap-backend/routes"
)

// @title MyCap API
// @version 1.0
// @description This is an API for MyCap Application

// @contact.name Dino Puguh
// @contact.email dinopuguh@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	err := database.Connect()
	if err != nil {
		panic("Can't connect database.")
	}

	migrate := flag.Bool("migrate", false, "Migrate all models to database")
	flag.Parse()

	if *migrate {
		migrations.All()
	}

	port := os.Getenv("PORT")

	app := routes.New()
	log.Fatal(app.Listen(":" + port))
}
