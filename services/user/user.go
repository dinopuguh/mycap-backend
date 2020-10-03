package user

import (
	"fmt"
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
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} response.HTTP{data=[]User}
// @Router /v1/users [get]
func GetAll(c *fiber.Ctx) error {
	db := database.DBConn

	var users []User
	if res := db.Find(&users); res.Error != nil {
		return c.JSON(response.HTTP{
			Status:  http.StatusServiceUnavailable,
			Message: res.Error.Error(),
		})
	}

	return c.JSON(response.HTTP{
		Success: true,
		Data:    users,
		Status:  http.StatusOK,
		Message: "Success get all users.",
	})
}

// Update function edit an user by ID
// @Summary Update user by ID
// @Description Update user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param user body UpdateUser true "Update user"
// @Success 200 {object} response.HTTP{data=User}
// @Router /v1/users [put]
func Update(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn

	var user User
	if res := db.First(&user, id); res.RowsAffected == 0 {
		return c.JSON(response.HTTP{
			Status:  http.StatusNotFound,
			Message: fmt.Sprintf("User with ID %v not found.", id),
		})
	}

	updatedUser := new(UpdateUser)
	if err := c.BodyParser(&updatedUser); err != nil {
		return c.JSON(response.HTTP{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	user.Name = updatedUser.Name
	user.RemainingTime = updatedUser.RemainingTime
	user.ReachedTimeLimit = updatedUser.ReachedTimeLimit

	db.Save(&user)

	return c.JSON(response.HTTP{
		Success: true,
		Data:    user,
		Status:  http.StatusOK,
		Message: "Success update user.",
	})
}

// Delete function removes an user by ID
// @Summary Remove user by ID
// @Description Remove user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.HTTP
// @Router /v1/users/{id} [delete]
func Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn

	var user User
	if res := db.First(&user, id); res.RowsAffected == 0 {
		return c.JSON(response.HTTP{
			Status:  http.StatusNotFound,
			Message: fmt.Sprintf("User with ID %v not found.", id),
		})
	}

	if res := db.Delete(&user); res.Error != nil {
		return c.JSON(response.HTTP{
			Status:  http.StatusServiceUnavailable,
			Message: res.Error.Error(),
		})
	}

	return c.JSON(response.HTTP{
		Success: true,
		Status:  http.StatusOK,
		Message: "Success delete user.",
	})
}
