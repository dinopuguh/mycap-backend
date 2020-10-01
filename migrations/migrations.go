package migrations

import (
	"log"

	"github.com/dinopuguh/mycap-backend/database"
	"github.com/dinopuguh/mycap-backend/services/user"
)

// All migrates all models to database
func All() {
	database.DBConn.AutoMigrate(&user.User{})

	log.Println("Models migrated to database.")
}
