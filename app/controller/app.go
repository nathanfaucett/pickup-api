package controller

import (
	"net/http"

	"github.com/aicacia/pickup/app"
	"github.com/aicacia/pickup/app/model"
	"github.com/aicacia/pickup/app/repository"
	"github.com/aicacia/pickup/docs"
	"github.com/gofiber/fiber/v2"
)

func GetOpenAPI(c *fiber.Ctx) error {
	c.Status(http.StatusOK)
	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.SendString(docs.SwaggerInfo.ReadDoc())
}

// GetHealthCheck
//
//	@Summary		Get Health Check
//	@Tags			app
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	model.HealthST
//	@Failure		500	{object}	model.HealthST
//	@Router			/health [get]
func GetHealthCheck(c *fiber.Ctx) error {
	health := model.HealthST{
		DB: repository.ValidConnection(),
	}
	if health.IsHealthy() {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusInternalServerError)
	}
	return c.JSON(health)
}

// GetVersion
//
//	@Summary		Get Version
//	@Tags			app
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	model.VersionST
//	@Router			/version [get]
func GetVersion(c *fiber.Ctx) error {
	c.Status(http.StatusOK)
	return c.JSON(app.Version)
}
