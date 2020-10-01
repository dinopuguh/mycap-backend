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
// @Success 200 {object} ResponseAuth
// @Failure 400 {object} response.HTTPError
// @Failure 503 {object} response.HTTPError
// @Router /v1/register [post]
func New(c *fiber.Ctx) error {
	db := database.DBConn

	user := new(User)
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.HTTPError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	var currentUser User
	if res := db.Where("email = ?", user.Email).First(&currentUser); res.RowsAffected > 0 {
		return c.Status(http.StatusBadRequest).JSON(response.HTTPError{
			Status:  http.StatusBadRequest,
			Message: "User with this email is already exist.",
		})
	}

	if res := db.Where("username = ?", user.Username).First(&currentUser); res.RowsAffected > 0 {
		return c.Status(http.StatusBadRequest).JSON(response.HTTPError{
			Status:  http.StatusBadRequest,
			Message: "User with this username is already exist.",
		})
	}

	var err error
	user.Password, err = helpers.HashPassword(user.Password)
	if err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(response.HTTPError{
			Status:  http.StatusServiceUnavailable,
			Message: err.Error(),
		})
	}

	if res := db.Create(user); res.Error != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(response.HTTPError{
			Status:  http.StatusServiceUnavailable,
			Message: res.Error.Error(),
		})
	}

	token, err := auth.GenerateJWT(user.Name, user.Email)
	if err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(response.HTTPError{
			Status:  http.StatusServiceUnavailable,
			Message: err.Error(),
		})
	}

	return c.JSON(ResponseAuth{
		User:        *user,
		AccessToken: token,
	})
}

// Login signs user to a session
// @Summary User login
// @Description User login
// @Tags auth
// @Accept json
// @Produce json
// @Param user body LoginUser true "User login"
// @Success 200 {object} ResponseAuth
// @Failure 400 {object} response.HTTPError
// @Failure 401 {object} response.HTTPError
// @Failure 404 {object} response.HTTPError
// @Failure 503 {object} response.HTTPError
// @Router /v1/login [post]
func Login(c *fiber.Ctx) error {
	db := database.DBConn

	login := new(LoginUser)
	if err := c.BodyParser(&login); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.HTTPError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	var user User
	if res := db.Where("email = ?", login.Email).First(&user); res.RowsAffected == 0 {
		return c.Status(http.StatusNotFound).JSON(response.HTTPError{
			Status:  http.StatusNotFound,
			Message: "User with this email not found.",
		})
	}

	if !helpers.CheckPasswordHash(login.Password, user.Password) {
		return c.Status(http.StatusUnauthorized).JSON(response.HTTPError{
			Status:  http.StatusUnauthorized,
			Message: "Password incorrect.",
		})
	}

	token, err := auth.GenerateJWT(user.Name, user.Email)
	if err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(response.HTTPError{
			Status:  http.StatusServiceUnavailable,
			Message: err.Error(),
		})
	}

	return c.JSON(ResponseAuth{
		User:        user,
		AccessToken: token,
	})
}
