package main

// @title Sample Cooking Recipes API
// @version 1.0
// @description  Support Vietnamese and English via locale Query ["vi" , "en"]
//
// @contact.name khanhanhtr
// @contact.url khanhanh.me
// @contact.email khanhanhtr00@gmail.com
//
// @schemes http https
// @accept json
// @accept json
import (
	"fmt"
	"khanhanhtr/sample/api/route"
	"khanhanhtr/sample/config"
	"khanhanhtr/sample/model"
	"khanhanhtr/sample/mongo"
	"khanhanhtr/sample/redis"
	"khanhanhtr/sample/repository"
	"khanhanhtr/sample/translator"
	"log"
	"os"
	"os/signal"
	"sync"

	logStack "github.com/diontr00/logstack"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

var (
	env_    *config.Env
	envOnce sync.Once

	mongo_    mongo.Client
	mongoOnce sync.Once

	redis_    redis.Client
	redisOnce sync.Once

	logger_    *logStack.Logger
	loggerOnce sync.Once

	translator_    *translator.UniversalTrans
	translatorOnce sync.Once

	validator_    *validator.Validate
	validatorOnce sync.Once

	userRepo_ model.UserRepository
	userOnce  sync.Once

	fiber_    *fiber.App
	fiberOnce sync.Once
)

func main() {
	envOnce.Do(func() {
		env_ = config.NewEnv()
	})

	var close_logger func() error

	loggerOnce.Do(func() {
		logger_, close_logger = config.NewLogger(env_)
	})

	mongoOnce.Do(func() {
		mongo_ = config.NewMongo(env_)
	})

	redisOnce.Do(func() {
		redis_ = config.NewRedis(env_)
	})

	translatorOnce.Do(func() {
		translator_ = config.NewTrans()
	})

	validatorOnce.Do(func() {
		validator_ = config.NewValidator()
	})

	fiberOnce.Do(func() {
		fiber_ = config.NewFiber(logger_, translator_)
	})

	userOnce.Do(func() {
		userRepo_ = repository.NewAuth0UserRepo(&env_.Auth0)
	})

	deps := &config.AppDependencies{
		Env:        env_,
		Mongo:      mongo_,
		Redis:      redis_,
		Logger:     logger_,
		Translator: translator_,
		Fiber:      fiber_,
		Validator:  validator_,
	}

	app := config.NewApp(deps)

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	go func() {
		<-terminate
		fmt.Println("Gracefully shutdown , Doing Cleanup task.... 🧹")
		if err := fiber_.Shutdown(); err != nil {
			log.Fatal(err.Error())
		}
		err := app.Close(close_logger)
		if err != nil {
			log.Fatal(err)
		}

	}()

	route.SetUp(&route.RouterOptions{
		Env:        app.Env,
		Timeout:    env_.App.ContextTimeout,
		RecipeDb:   mongo_.UseDatabase(env_.Mongo.DBName),
		UserDb:     userRepo_,
		Logger:     logger_,
		Fiber:      fiber_,
		Validator:  validator_,
		Translator: translator_,
		Cache:      redis_,
	})

	err := fiber_.Listen(":8080")
	if err != nil {
		logger_.Fatal("[App]", logStack.String("error", err.Error()))
	}

}
