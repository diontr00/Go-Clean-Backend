package controller

import (
	"fmt"
	"khanhanhtr/sample/internal/tokenutil"
	"khanhanhtr/sample/model"
	"khanhanhtr/sample/redis"
	"khanhanhtr/sample/translator"
	"time"

	logStack "github.com/diontr00/logstack"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type RefreshController struct {
	Logger        *logStack.Logger
	Validator     *validator.Validate
	Translator    translator.Translator
	RecipeUseCase model.RecipeUseCase
	Cache         redis.Client
	EncryptKey    string
}

func (r *RefreshController) RefreshToken(c *fiber.Ctx) error {
	token, ok := c.Locals("userClaim").(*jwt.Token)
	if !ok {

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": r.Translator.TranslateMessage(c, "missingtoken", nil, nil),
		})
	}
	claims, ok := token.Claims.(*model.JWTCustomClaims)

	if !ok {

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": r.Translator.TranslateMessage(c, "missingtoken", nil, nil),
		})
	}

	username := claims.Username

	expiry, err := claims.GetExpirationTime()
	fmt.Println(expiry)
	if err != nil || expiry.Unix() > time.Now().Add(5*time.Minute).Unix() {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": r.Translator.TranslateMessage(c, "tryagain", nil, nil),
		})

	}

	refresh_token, err := tokenutil.CreateAccessToken(
		username,
		r.EncryptKey,
		time.Now().Add(time.Minute*10),
	)

	if err != nil {

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": r.Translator.TranslateMessage(c, "tryagain", nil, nil),
		})

	}

	return c.JSON(refresh_token)
}
