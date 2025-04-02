package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/kooroshh/fiber-boostrap/app/controllers"
)

type ApiRouter struct{}

func (h ApiRouter) InstallRouter(app *fiber.App) {
	api := app.Group("/api", limiter.New())
	api.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Hello from api",
		})
	})

	userGroup := app.Group("/user")
	userV1Group := userGroup.Group("/v1")
	userV1Group.Post("/register", controllers.Register)
	userV1Group.Post("/login", controllers.Login)
	userV1Group.Delete("/logout", MiddlewareValidateAuth, controllers.Logout)
}

func NewApiRouter() *ApiRouter {
	return &ApiRouter{}
}
