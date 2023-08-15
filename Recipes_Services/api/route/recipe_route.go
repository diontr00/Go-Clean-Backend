package route

import (
	"khanhanhtr/sample/api/controller"
	"khanhanhtr/sample/mongo"
	"khanhanhtr/sample/redis"
	"khanhanhtr/sample/repository"
	"khanhanhtr/sample/translator"
	"khanhanhtr/sample/usecase"
	"time"

	logStack "github.com/diontr00/logstack"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type RecipeRouteConfig struct {
	Db         mongo.Database
	Logger     *logStack.Logger
	Validator  *validator.Validate
	Translator translator.Translator
	Group      fiber.Router
	Cache      redis.Client
	EncryptKey string
	Timeout    time.Duration
}

func recipeSetup(
	config *RecipeRouteConfig,
) {
	repo := repository.NewMongoRecipeRepo(config.Db, "recipes")
	controller := &controller.RecipeController{

		RecipeUseCase: usecase.NewRecipeUsecase(repo, config.Logger, config.Timeout, config.Cache),
		Validator:     config.Validator,
		Translator:    config.Translator,
		Logger:        config.Logger,
	}

	config.Group.Get("/", controller.Get)
	config.Group.Get("/item/:id?", controller.GetById)
	config.Group.Post("/", controller.Create)
	config.Group.Put("/item/:id?", controller.Update)
	config.Group.Delete("/item/:id?", controller.Delete)

}
