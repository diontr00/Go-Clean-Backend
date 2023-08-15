package middleware

import (
	"crypto/hmac"
	"khanhanhtr/sample/internal/jwk"
	"khanhanhtr/sample/internal/tokenutil"
	"khanhanhtr/sample/model"
	"khanhanhtr/sample/translator"

	logStack "github.com/diontr00/logstack"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func NewJWTMiddleWare(
	logger *logStack.Logger,
	translator translator.Translator,
	provider jwk.JWKSProvider,

) func(c *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		Claims:      &model.JWTCustomClaims{},
		ContextKey:  "user-info",
		TokenLookup: "cookie:Authorization",
		KeyFunc:     provider.KeyFunc(jwk.RS256),
		Filter: func(c *fiber.Ctx) bool {

			if c.Cookies("Authorization") == "" {
				c.Locals("authenticated", false)
				return true
			}
			return false
			// 	supportedPaths := map[string]func(c *fiber.Ctx) bool{
			// 		"/login": func(c *fiber.Ctx) bool {
			// 			return c.Cookies("Authorization") == ""
			// 		},
			// 	}
			//
			// 	path := c.Request().URI().Path()
			// 	skipCondition, ok := supportedPaths[string(path)]
			// 	if !ok {
			// 		return false
			// 	}
			//
			// 	return skipCondition(c)
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.Error(err.Error())
			c.Locals("authenticated", false)
			return c.Next()
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			c.Locals("authenticated", true)
			return c.Next()
		},
	},
	)
}

func NewSessionMiddleWare(
	sessionStore *session.Store,
	logger *logStack.Logger,
	translator translator.Translator,
	secret []byte,
) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sessionId := c.Cookies("session_id")
		if sessionId == "" {
			c.Locals("authenticated", false)
			return c.Next()
		}
		session, err := sessionStore.Get(c)
		if err != nil {
			logger.Error(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": translator.TranslateMessage(c, "tryagain", nil, nil),
			})
		}

		if session.Fresh() || session.Get("Access_Token") == nil {
			c.Locals("authenticated", false)
			return c.Next()
		}

		calculatedSignature := tokenutil.GenerateSignature(sessionId, secret)
		storedSignature := session.Get("session_signature")

		if !hmac.Equal([]byte(storedSignature.(string)), []byte(calculatedSignature)) {
			c.Locals("authenticated", false)
			return c.Next()
		}

		c.Locals("authenticated", true)
		return c.Next()
	}
}
