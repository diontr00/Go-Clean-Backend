package route

import (
	"khanhanhtr/sample2/api/middleware"
	"khanhanhtr/sample2/config"
	mqueue "khanhanhtr/sample2/rabbitmq"
	"khanhanhtr/sample2/translator"
	"time"

	logStack "github.com/diontr00/logstack"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type RouterOptions struct {
	Env        *config.Env
	Timeout    time.Duration
	Logger     *logStack.Logger
	Validator  *validator.Validate
	Translator translator.Translator
	Fiber      *fiber.App
	Mq         mqueue.Mqclient
}

func Setup(setup *RouterOptions) {
	// Recover
	setup.Fiber.Use(recover.New(
		recover.Config{
			EnableStackTrace: true,
			StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
				setup.Logger.Error("Panic Recover", logStack.Any("stack:", e))
			},
		},
	))

	// Locale
	setup.Fiber.Use(middleware.NewLocaleMiddleWare)

	rssGroup := setup.Fiber.Group("/rss")
	rssSetup(

		&RssRouteConfig{
			Logger:     setup.Logger,
			Validator:  setup.Validator,
			Translator: setup.Translator,
			Group:      rssGroup,
			Timeout:    setup.Env.App.ContextTimeout,
			Mq:         setup.Mq,
		},
	)

	setup.Fiber.All("*", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"error": setup.Translator.TranslateMessage(c, "notsupport", translator.ParamConfig{
				"Method": c.Method()}, nil),
		})
	})

}
