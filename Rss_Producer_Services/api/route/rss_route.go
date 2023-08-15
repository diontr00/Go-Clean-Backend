package route

import (
	"khanhanhtr/sample2/api/controller"
	mqueue "khanhanhtr/sample2/rabbitmq"
	"khanhanhtr/sample2/translator"
	"khanhanhtr/sample2/usecase"
	"time"

	logStack "github.com/diontr00/logstack"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type RssRouteConfig struct {
	Logger     *logStack.Logger
	Validator  *validator.Validate
	Translator translator.Translator
	Group      fiber.Router
	Timeout    time.Duration
	Mq         mqueue.Mqclient
}

func rssSetup(
	config *RssRouteConfig,
) {

	controller := &controller.RssController{

		RssUsecase: usecase.NewRssUsecase(config.Logger, config.Mq),
		Validator:  config.Validator,
		Translator: config.Translator,
		Logger:     config.Logger,
	}

	config.Group.Post("/fetch", controller.ParseRss)

}
