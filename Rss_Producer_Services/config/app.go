package config

import (
	mqueue "khanhanhtr/sample2/rabbitmq"
	"khanhanhtr/sample2/translator"
	"log"
	"os"

	logStack "github.com/diontr00/logstack"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-multierror"
)

type AppDependencies struct {
	Env        *Env
	Logger     *logStack.Logger
	Translator translator.Translator
	Validator  *validator.Validate
	Messqueue  mqueue.Mqclient
	Fiber      *fiber.App
}

type applications struct {
	Env        *Env
	Logger     *logStack.Logger
	Translator translator.Translator
	Validator  *validator.Validate
	Messqueue  mqueue.Mqclient
	Fiber      *fiber.App
}

func NewApp(deps *AppDependencies) *applications {
	app := &applications{Env: deps.Env,
		Logger:     deps.Logger,
		Translator: deps.Translator,
		Validator:  deps.Validator,
		Messqueue:  deps.Messqueue,
	}
	return app
}

// Setup logger
// ---------------------------------------------------------------------------------------------------------------

func NewLogger(env *Env) (*logStack.Logger, func() error) {
	var logger *logStack.Logger
	if env.App.LogLocation == "" {
		logger = logStack.DefaultLogger()
		return logger, os.Stdout.Close
	} else {
		log_file, err := os.OpenFile(env.App.LogLocation, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			logger = logStack.DefaultLogger()
			logger.Fatal("[LOG FILE]", logStack.String("error:", err.Error()))
		}

		logger = logStack.NewLogger(log_file, logStack.InfoLevel)

		return logger, log_file.Close
	}
}

// Setup Translator and Validator
// -------------------------------------------------------------------------------------------------------------------

func NewTrans() *translator.UniversalTrans {

	translator, err := newTranslator()
	if err != nil {
		log.Fatal("[TRANSLATOR]", logStack.String("error", err.Error()))
	}
	return translator

}

func NewValidator() *validator.Validate {
	return validator.New()
}

func NewFiber(logger *logStack.Logger, translator translator.Translator) *fiber.App {

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.Error(err.Error())
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": translator.TranslateMessage(c, "notfound", nil, nil),
			})

		},
	})
	return app
}

// Setup RabbitMQ
// -----------------------------------------------------------------------------------
func NewRabbitMQ(env *Env) mqueue.Mqclient {

	return *newMessageQueue(env)
}

// Close the application by closing logfile and mongo connection
func (a applications) Close(logger func() error) error {
	var errs error

	if err := logger(); err != nil {
		errs = multierror.Append(errs, err)
	}

	if err := a.Messqueue.Close(); err != nil {

		errs = multierror.Append(errs, err)

	}

	_, ok := errs.(*multierror.Error)
	if !ok {
		errs = nil
		return errs
	}

	return errs.(*multierror.Error).ErrorOrNil()
}
