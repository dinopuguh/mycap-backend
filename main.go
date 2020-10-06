package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/dinopuguh/mycap-backend/database"
	_ "github.com/dinopuguh/mycap-backend/docs"
	"github.com/dinopuguh/mycap-backend/migrations"
	"github.com/dinopuguh/mycap-backend/routes"
	"github.com/dinopuguh/mycap-backend/seed"
	"github.com/dinopuguh/mycap-backend/services/user"
	"github.com/go-co-op/gocron"
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
	if err := database.Connect(); err != nil {
		panic("Can't connect database.")
	}

	migrate := flag.Bool("migrate", false, "Migrate all models to database")
	flag.Parse()

	if *migrate {
		migrations.All()
	}

	for _, seeder := range seed.AllTypes() {
		if err := seeder.Run(database.DBConn); err != nil {
			log.Fatalln(err.Error())
		}
	}

	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(1).Month(7).Do(user.ResetTimeLimit)
	scheduler.StartAsync()

	port := os.Getenv("PORT")
	app := routes.New()
	log.Fatal(app.Listen(":" + port))
}
