package seed

import "gorm.io/gorm"

// Seed is a struct to handle seeder
type Seed struct {
	Name string
	Run  func(*gorm.DB) error
}
