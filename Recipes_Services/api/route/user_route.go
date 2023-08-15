package route

import (
	"khanhanhtr/sample/api/controller"
	"khanhanhtr/sample/config"
	"khanhanhtr/sample/model"
	"khanhanhtr/sample/redis"
	"khanhanhtr/sample/translator"
	"khanhanhtr/sample/usecase"
	"time"

	logStack "github.com/diontr00/logstack"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type UserRouteConfig struct {
	Db         model.UserRepository
	Logger     *logStack.Logger
	Validator  *validator.Validate
	Translator translator.Translator
	Group      fiber.Router
	Cache      redis.Client
	EncryptKey string
	Timeout    time.Duration
	Auth       *config.Auth0Env
}

func userSetup(
	config *UserRouteConfig,
) {
	controller := &controller.UserController{

		Logger:      config.Logger,
		Validator:   config.Validator,
		Translator:  config.Translator,
		UserUseCase: usecase.NewUserUseCase(config.Timeout, config.Db),
		EncryptKey:  config.EncryptKey,
		Signature:   []byte(config.EncryptKey),
	}

	config.Group.Post("/signin", controller.Signin)
	config.Group.Post("/signout", controller.SignOut)
	config.Group.Post("/signup", controller.Signup)
}
