package scheduler

import (
	"log"
	"os/user"

	"github.com/dinopuguh/mycap-backend/database"
)

// ResetTimeLimit function updates users' time limit monthly
func ResetTimeLimit() {
	db := database.DBConn

	db.Model(user.User{}).Where("type_id = ?", 1).Updates(map[string]interface{}{
		"reached_time_limit": false,
		"remaining_time":     36000000,
	})
	log.Println("Update free users' remaining time.")
}
