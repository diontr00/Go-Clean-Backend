package config

import (
	"context"
	"khanhanhtr/sample2/api/controller"
	"khanhanhtr/sample2/mongo"
	mqueue "khanhanhtr/sample2/rabbitmq"
	"khanhanhtr/sample2/repository"
	"log"
	"os"

	logStack "github.com/diontr00/logstack"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-multierror"
)

type AppDependencies struct {
	Env       *Env
	Mongo     mongo.Client
	Logger    *logStack.Logger
	Messqueue mqueue.Mqclient
	Fiber     *fiber.App
}

type applications struct {
	Env       *Env
	Mongo     mongo.Client
	Logger    *logStack.Logger
	Messqueue mqueue.Mqclient
	Fiber     *fiber.App
}

func NewApp(deps *AppDependencies) *applications {
	app := &applications{Env: deps.Env,
		Mongo:     deps.Mongo,
		Logger:    deps.Logger,
		Messqueue: deps.Messqueue,
	}
	return app
}

func (a *applications) StartRss(done chan struct{}) {
	log.Println("Running Rss Services...🚀")
	rss_repo := repository.NewMongoRssRepo(a.Mongo.UseDatabase("restapi"), "rss-feed")

	controller := controller.NewRssController(
		rss_repo,
		a.Logger,
		a.Env.App.ContextTimeout,
		a.Messqueue,
	)

	controller.ListenForRssRequest(context.Background(), done)

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

// Setup DB
// ----------------------------------------------------------------------------------------------------------------

func NewMongo(env *Env) mongo.Client {
	mongodb, err := NewMongoDatabase(env.Mongo)
	if err != nil {
		log.Fatalf("Error initialize mongo db: %v", err)
	}
	return mongodb
}

// Setup RabbitMQ
// -----------------------------------------------------------------------------------
func NewRabbitMQ(env *Env) mqueue.Mqclient {

	return *newMessageQueue(env)
}

// Close the application by closing logfile and mongo connection
func (a applications) Close(logger func() error) error {
	var errs error

	if a.Mongo != nil {
		err := a.Mongo.Disconnect(context.Background())
		errs = multierror.Append(errs, err)
	}

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
