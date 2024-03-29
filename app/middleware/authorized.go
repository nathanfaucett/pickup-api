package middleware

import (
	"log"
	"net/http"

	"github.com/aicacia/pickup/app/jwt"
	"github.com/aicacia/pickup/app/model"
	"github.com/aicacia/pickup/app/repository"
	"github.com/gofiber/fiber/v2"
)

var claimsLocalKey = "claims"
var userLocalKey = "user"

func AuthorizedMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := GetAuthorizationFromContext(c)
		claims, err := jwt.ParseClaimsFromToken(tokenString)
		if err != nil {
			log.Printf("failed to parse claims from token: %v", err)
			return model.NewError(http.StatusUnauthorized).AddError("authorization", "invalid", "header").Send(c)
		}
		user, err := repository.GetUserById(claims.Subject)
		if err != nil {
			log.Printf("failed to fetch user: %v", err)
			return model.NewError(http.StatusUnauthorized).AddError("authorization", "invalid", "header").Send(c)
		}
		c.Locals(userLocalKey, user)
		c.Locals(claimsLocalKey, claims)
		return c.Next()
	}
}

func GetClaims(c *fiber.Ctx) *jwt.Claims {
	claims := c.Locals(claimsLocalKey)
	return claims.(*jwt.Claims)
}

func GetUser(c *fiber.Ctx) *repository.UserRowST {
	user := c.Locals(userLocalKey)
	return user.(*repository.UserRowST)
}
