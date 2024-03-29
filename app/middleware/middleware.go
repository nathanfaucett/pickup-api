package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

var bearerString = "Bearer "

func GetAuthorizationFromContext(c *fiber.Ctx) string {
	authorizationHeader := strings.TrimSpace(c.Get("Authorization"))
	if len(authorizationHeader) != 0 {
		return strings.TrimSpace(authorizationHeader[len(bearerString):])
	} else {
		return ""
	}
}
