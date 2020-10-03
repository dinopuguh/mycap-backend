package migrations

import (
	"log"

	"github.com/dinopuguh/mycap-backend/database"
	"github.com/dinopuguh/mycap-backend/services/group"
	"github.com/dinopuguh/mycap-backend/services/user"
)

// All migrates all models to database
func All() {
	database.DBConn.AutoMigrate(&user.User{})
	database.DBConn.AutoMigrate(&group.Group{})

	log.Println("Models migrated to database.")
}
