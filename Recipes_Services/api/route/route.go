package route

import (
	"khanhanhtr/sample/api/middleware"
	"khanhanhtr/sample/config"
	"khanhanhtr/sample/internal/jwk"
	"khanhanhtr/sample/model"
	"khanhanhtr/sample/mongo"
	"khanhanhtr/sample/redis"
	"khanhanhtr/sample/translator"
	"time"

	logStack "github.com/diontr00/logstack"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type RouterOptions struct {
	Env        *config.Env
	Timeout    time.Duration
	RecipeDb   mongo.Database
	UserDb     model.UserRepository
	Logger     *logStack.Logger
	Validator  *validator.Validate
	Translator translator.Translator
	Fiber      *fiber.App
	Cache      redis.Client
}

func SetUp(
	setup *RouterOptions,

) {

	setup.Fiber.Use(recover.New(
		recover.Config{
			EnableStackTrace: true,
			StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
				setup.Logger.Error("Panic Recover", logStack.Any("stack:", e))
			},
		},
	))

	setup.Fiber.Use(middleware.NewLocaleMiddleWare)

	// session_store := session.New(session.Config{
	// 	Expiration:     10 * time.Minute,
	// 	KeyLookup:      "cookie:session_id",
	// 	KeyGenerator:   random.RandomStringMust,
	// 	CookieSecure:   true,
	// 	CookieHTTPOnly: true,
	// 	CookieSameSite: "Strict",
	// })

	setup.Fiber.Use(
		middleware.NewCookieEncryptMiddleWare(setup.Env.App.CookieEncryptKey, setup.Logger),
	)
	// setup.Fiber.Use(
	// 	middleware.NewSessionMiddleWare(
	// 		session_store,
	// 		setup.Logger,
	// 		setup.Translator,
	// 		[]byte(setup.Env.App.EncryptKey),
	// 	),
	// )

	jwks_url, err := jwk.NewCommonURL(setup.Env.Auth0.Domain)
	if err != nil {
		setup.Logger.Fatal(err.Error())
	}
	jwks_provider := jwk.NewJWKProviderWithCache(jwks_url, time.Minute*5)

	setup.Fiber.Use(
		middleware.NewJWTMiddleWare(setup.Logger, setup.Translator, jwks_provider),
	)

	// setup.Fiber.Use(csrf.New(csrf.Config{
	// 	KeyLookup:      "header:X-Csrf-Token",
	// 	CookieName:     "csrf_",
	// 	CookieSameSite: "Strict",
	// 	Expiration:     3 * time.Hour,
	// 	KeyGenerator:   random.RandomStringMust,
	// }))

	userGroup := setup.Fiber.Group("/auth")
	recipeGroup := setup.Fiber.Group("/recipes")
	recipeSetup(
		&RecipeRouteConfig{
			Db:         setup.RecipeDb,
			Logger:     setup.Logger,
			Validator:  setup.Validator,
			Translator: setup.Translator,
			Group:      recipeGroup,
			Cache:      setup.Cache,
			Timeout:    setup.Env.App.ContextTimeout,
		},
	)

	userSetup(
		&UserRouteConfig{
			Db:         setup.UserDb,
			Logger:     setup.Logger,
			Validator:  setup.Validator,
			Cache:      setup.Cache,
			Translator: setup.Translator,
			Group:      userGroup,
			EncryptKey: setup.Env.App.EncryptKey,
			Timeout:    setup.Env.App.ContextTimeout,
		},
	)

	setup.Fiber.All("*", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"error": setup.Translator.TranslateMessage(c, "notsupport", translator.ParamConfig{
				"Method": c.Method()}, nil),
		})
	})

}
