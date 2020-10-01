package user

import (
	"net/http"

	"github.com/dinopuguh/mycap-backend/database"
	"github.com/dinopuguh/mycap-backend/response"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// User is a model for user
type User struct {
	gorm.Model
	Name             string `json:"name"`
	Username         string `json:"username"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	RemainingTime    int64  `json:"remaining_time" gorm:"default:18000;"`
	ReachedTimeLimit bool   `json:"reached_time_limit" gorm:"default:false;"`
}

// GetAll is a function to get all users data from database
// @Summary Get all users
// @Description Get all users
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} []User
// @Failure 503 {object} response.HTTPError
// @Router /v1/users [get]
func GetAll(c *fiber.Ctx) error {
	db := database.DBConn

	var users []User
	if res := db.Find(&users); res.Error != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(response.HTTPError{
			Status:  http.StatusServiceUnavailable,
			Message: res.Error.Error(),
		})
	}

	return c.JSON(users)
}
