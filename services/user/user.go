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
	RemainingTime    int64  `json:"remaining_time" gorm:"default:36000000;"`
	ReachedTimeLimit bool   `json:"reached_time_limit" gorm:"default:false;"`
	Type             Type   `json:"type"`
	TypeID           uint   `json:"type_id"`
}

// Type is a model for user's type
type Type struct {
	gorm.Model
	Name string `json:"name"`
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
	if res := db.Preload("Type").Find(&users); res.Error != nil {
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
// @Param id path int true "User ID"
// @Param user body UpdateUser true "Update user"
// @Success 200 {object} response.HTTP{data=User}
// @Security ApiKeyAuth
// @Router /v1/users/{id} [put]
func Update(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn

	updatedUser := new(UpdateUser)
	if err := c.BodyParser(&updatedUser); err != nil {
		return c.JSON(response.HTTP{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	user := new(User)
	if err := db.Preload("Type").First(&user, id).Error; err != nil {
		switch err.Error() {
		case "record not found":
			return c.JSON(response.HTTP{
				Status:  http.StatusNotFound,
				Message: fmt.Sprintf("User with ID %v not found.", id),
			})
		default:
			return c.JSON(response.HTTP{
				Status:  http.StatusServiceUnavailable,
				Message: err.Error(),
			})
		}
	}

	userType := new(Type)
	if updatedUser.TypeID != 0 {
		if res := db.First(&userType, updatedUser.TypeID); res.RowsAffected == 0 {
			return c.JSON(response.HTTP{
				Status:  http.StatusNotFound,
				Message: fmt.Sprintf("User type with ID %v not found.", updatedUser.TypeID),
			})
		}

		user.TypeID = updatedUser.TypeID
		user.Type = *userType
	}

	user.Name = updatedUser.Name
	user.ReachedTimeLimit = updatedUser.ReachedTimeLimit
	user.RemainingTime = updatedUser.RemainingTime

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
// @Security ApiKeyAuth
// @Router /v1/users/{id} [delete]
func Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn

	var user User
	if err := db.First(&user, id).Error; err != nil {
		switch err.Error() {
		case "record not found":
			return c.JSON(response.HTTP{
				Status:  http.StatusNotFound,
				Message: fmt.Sprintf("User with ID %v not found.", id),
			})
		default:
			return c.JSON(response.HTTP{
				Status:  http.StatusServiceUnavailable,
				Message: err.Error(),
			})

		}
	}

	db.Delete(&user)

	return c.JSON(response.HTTP{
		Success: true,
		Status:  http.StatusOK,
		Message: "Success delete user.",
	})
}
