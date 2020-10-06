package seed

import (
	"log"

	"github.com/dinopuguh/mycap-backend/services/user"
	"gorm.io/gorm"
)

func createType(db *gorm.DB, name string) error {
	userType := new(user.Type)
	if res := db.Where("name = ?", name).First(&userType); res.RowsAffected != 0 {
		log.Printf("Type %s is already exist.\n", name)
		return nil
	}

	return db.Create(&user.Type{Name: name}).Error
}

// AllTypes function return all type's seeds
func AllTypes() []Seed {
	return []Seed{
		Seed{
			Name: "Create user type Free",
			Run: func(db *gorm.DB) error {
				return createType(db, "Free")
			},
		},
		Seed{
			Name: "Create user type Premium",
			Run: func(db *gorm.DB) error {
				return createType(db, "Premium")
			},
		},
		Seed{
			Name: "Create user type Pro",
			Run: func(db *gorm.DB) error {
				return createType(db, "Pro")
			},
		},
	}
}
