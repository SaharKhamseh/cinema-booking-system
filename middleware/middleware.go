package middleware

import (
	"strconv"

	"github.com/SaharKhamseh/cinema-backend/database"
	"github.com/SaharKhamseh/cinema-backend/models"
	"github.com/SaharKhamseh/cinema-backend/util"
	"github.com/gofiber/fiber/v2"
)

func IsAuthentication(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	if _, err := util.Parsejwt(cookie); err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	return c.Next()
}

func IsAdmin(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	id, err := util.Parsejwt(cookie)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	userId, _ := strconv.Atoi(id)

	var user models.User
	database.DB.First(&user, userId)

	if user.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access Denied",
		})
	}
	return c.Next()
}
