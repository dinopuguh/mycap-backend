package routes

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/dinopuguh/mycap-backend/auth"
	"github.com/dinopuguh/mycap-backend/services/group"
	"github.com/dinopuguh/mycap-backend/services/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v2"
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

	v1.Get("/groups", group.GetAll)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: auth.SigningKey,
	}))

	v1.Put("/users/:id", user.Update)

	v1.Post("/groups", group.New)
	v1.Post("/join-groups", group.Join)
	v1.Post("/leave-groups", group.Leave)

	return app
}
