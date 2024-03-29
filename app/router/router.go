package router

import (
	"github.com/aicacia/pickup/app/config"
	"github.com/aicacia/pickup/app/controller"
	"github.com/aicacia/pickup/app/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func InstallRouter(fiberApp *fiber.App) {
	root := fiberApp.Group("", cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
		AllowCredentials: true,
	}))

	if config.Get().OpenAPI.Enabled {
		root.Get("/openapi.json", controller.GetOpenAPI)
	}

	root.Get("/health", controller.GetHealthCheck)
	root.Get("/version", controller.GetVersion)

	oauth2 := root.Group("/oauth2")
	oauth2.Get("/:provider", controller.GetProviderRedirect)
	oauth2.Get("/:provider/callback", controller.GetProviderCallback)

	authenticated := root.Group("")
	authenticated.Use(middleware.AuthorizedMiddleware())

	user := authenticated.Group("/user")
	user.Get("", controller.GetCurrentUser)
	user.Patch("", controller.PatchUpdateUser)
	user.Patch("/complete", controller.PatchCompleteUser)

	places := authenticated.Group("/places")
	places.Get("", controller.GetPlaces)
}
