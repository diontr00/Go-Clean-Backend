package middleware

import "github.com/gofiber/fiber/v2"

func NewLocaleMiddleWare(c *fiber.Ctx) error {
	locale := c.Query("locale", "en")
	switch locale {
	case "en":
		locale = "en-US"
	case "vi":
		locale = "vi-VN"
	default:
		locale = "en-US"
	}
	c.Locals("locale", locale)
	return c.Next()
}
