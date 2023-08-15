package middleware

import (
	logStack "github.com/diontr00/logstack"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
)

// Key must be  32 character string
func NewCookieEncryptMiddleWare(
	secret string,
	logger *logStack.Logger,
) func(c *fiber.Ctx) error {

	return encryptcookie.New(
		encryptcookie.Config{
			Key: secret,
		},
	)

}
