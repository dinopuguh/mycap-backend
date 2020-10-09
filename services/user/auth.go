package user

import (
	"net/http"

	"github.com/dinopuguh/mycap-backend/auth"
	"github.com/dinopuguh/mycap-backend/database"
	"github.com/dinopuguh/mycap-backend/helpers"
	"github.com/dinopuguh/mycap-backend/response"
	"github.com/gofiber/fiber/v2"
)

// ResponseAuth represents response body for authenticated user
type ResponseAuth struct {
	User        User   `json:"user"`
	AccessToken string `json:"access_token"`
}

// New registers a new user data
// @Summary Register a new user
// @Description Register user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterUser true "Register user"
// @Success 200 {object} response.HTTP{data=ResponseAuth}
// @Router /v1/register [post]
func New(c *fiber.Ctx) error {
	db := database.DBConn

	registerUser := new(RegisterUser)
	if err := c.BodyParser(&registerUser); err != nil {
		return c.JSON(response.HTTP{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if registerUser.TypeID == 0 {
		registerUser.TypeID = 1
	}

	userType := new(Type)
	if err := db.First(&userType, registerUser.TypeID).Error; err != nil {
		switch err.Error() {
		case "record not found":
			return c.JSON(response.HTTP{
				Status:  http.StatusBadRequest,
				Message: "User type with this ID not exist.",
			})
		default:
			return c.JSON(response.HTTP{
				Status:  http.StatusServiceUnavailable,
				Message: err.Error(),
			})
		}
	}

	user := new(User)
	if res := db.Where("email = ?", registerUser.Email).First(&user); res.RowsAffected > 0 {
		return c.JSON(response.HTTP{
			Status:  http.StatusBadRequest,
			Message: "User with this email is already exist.",
		})
	}

	if res := db.Where("username = ?", registerUser.Username).First(&user); res.RowsAffected > 0 {
		return c.JSON(response.HTTP{
			Status:  http.StatusBadRequest,
			Message: "User with this username is already exist.",
		})
	}

	user.Name = registerUser.Name
	user.Email = registerUser.Email
	user.Username = registerUser.Username
	user.Type = *userType

	var err error
	user.Password, err = helpers.HashPassword(registerUser.Password)
	if err != nil {
		return c.JSON(response.HTTP{
			Status:  http.StatusServiceUnavailable,
			Message: err.Error(),
		})
	}

	db.Create(user)

	token, err := auth.GenerateJWT(user.Name, user.Email)
	if err != nil {
		return c.JSON(response.HTTP{
			Status:  http.StatusServiceUnavailable,
			Message: err.Error(),
		})
	}

	return c.JSON(response.HTTP{
		Success: true,
		Data: ResponseAuth{
			User:        *user,
			AccessToken: token,
		},
		Status:  http.StatusOK,
		Message: "Success register.",
	})
}

// Login signs user to a session
// @Summary User login
// @Description User login
// @Tags auth
// @Accept json
// @Produce json
// @Param user body LoginUser true "User login"
// @Success 200 {object} response.HTTP{data=ResponseAuth}
// @Router /v1/login [post]
func Login(c *fiber.Ctx) error {
	db := database.DBConn

	login := new(LoginUser)
	if err := c.BodyParser(&login); err != nil {
		return c.JSON(response.HTTP{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	var user User
	if res := db.Preload("Type").Where("email = ?", login.Email).First(&user); res.RowsAffected == 0 {
		return c.JSON(response.HTTP{
			Status:  http.StatusNotFound,
			Message: "User with this email not found.",
		})
	}

	if !helpers.CheckPasswordHash(login.Password, user.Password) {
		return c.JSON(response.HTTP{
			Status:  http.StatusUnauthorized,
			Message: "Password incorrect.",
		})
	}

	token, err := auth.GenerateJWT(user.Name, user.Email)
	if err != nil {
		return c.JSON(response.HTTP{
			Status:  http.StatusServiceUnavailable,
			Message: err.Error(),
		})
	}

	return c.JSON(response.HTTP{
		Success: true,
		Data: ResponseAuth{
			User:        user,
			AccessToken: token,
		},
		Status:  http.StatusOK,
		Message: "Success login.",
	})
}
