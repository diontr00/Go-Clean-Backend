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
	"khanhanhtr/sample2/api/route"
	"khanhanhtr/sample2/config"
	mqueue "khanhanhtr/sample2/rabbitmq"
	"khanhanhtr/sample2/translator"
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

	logger_    *logStack.Logger
	loggerOnce sync.Once

	translator_    *translator.UniversalTrans
	translatorOnce sync.Once

	validator_    *validator.Validate
	validatorOnce sync.Once

	fiber_    *fiber.App
	fiberOnce sync.Once

	mqueue_ mqueue.Mqclient
	mqOnce  sync.Once
)

func main() {
	envOnce.Do(func() {
		env_ = config.NewEnv()
	})

	var close_logger func() error

	loggerOnce.Do(func() {
		logger_, close_logger = config.NewLogger(env_)
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

	mqOnce.Do(func() {
		mqueue_ = config.NewRabbitMQ(env_)
	})

	deps := &config.AppDependencies{
		Env:        env_,
		Logger:     logger_,
		Translator: translator_,
		Fiber:      fiber_,
		Validator:  validator_,
		Messqueue:  mqueue_,
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

	route.Setup(&route.RouterOptions{
		Env:        app.Env,
		Timeout:    env_.App.ContextTimeout,
		Logger:     logger_,
		Fiber:      fiber_,
		Validator:  validator_,
		Mq:         mqueue_,
		Translator: translator_,
	})

	err := fiber_.Listen(":8080")
	if err != nil {
		logger_.Fatal("[App]", logStack.String("error", err.Error()))
	}

}
