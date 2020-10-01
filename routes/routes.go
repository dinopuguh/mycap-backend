package routes

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/dinopuguh/mycap-backend/services/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// New create an instance of MyCap routes
func New() *fiber.App {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New(logger.Config{
		Format:     "${cyan}[${time}] ${white}${pid} ${red}${status} ${blue}[${method}] ${white}${path}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "Asia/Jakarta",
	}))
	app.Use("/swagger", swagger.Handler)

	api := app.Group("/api")
	v1 := api.Group("/v1", func(c *fiber.Ctx) error {
		c.JSON(fiber.Map{
			"message": "🐣 v1",
		})
		return c.Next()
	})
	v1.Post("/register", user.New)
	v1.Post("/login", user.Login)

	v1.Get("/users", user.GetAll)
	v1.Delete("/users/:id", user.Delete)

	// app.Use(jwtware.New(jwtware.Config{
	// 	SigningKey: auth.MySigningKey,
	// 	Claims:     &auth.JwtCustomClaims{},
	// }))

	// app.Post("/api/v1/address", address.New)

	return app
}
