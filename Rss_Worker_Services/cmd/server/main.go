package main

import (
	"fmt"
	"khanhanhtr/sample2/config"
	"khanhanhtr/sample2/mongo"
	mqueue "khanhanhtr/sample2/rabbitmq"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/diontr00/logstack"
)

var (
	env_    *config.Env
	envOnce sync.Once

	mongo_    mongo.Client
	mongoOnce sync.Once

	logger_    *logStack.Logger
	loggerOnce sync.Once

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

	mongoOnce.Do(func() {
		mongo_ = config.NewMongo(env_)
	})

	mqOnce.Do(func() {
		mqueue_ = config.NewRabbitMQ(env_)
	})

	deps := &config.AppDependencies{
		Env:       env_,
		Mongo:     mongo_,
		Logger:    logger_,
		Messqueue: mqueue_,
	}

	app := config.NewApp(deps)

	doneCh := make(chan struct{}, 1)

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	go func() {
		<-terminate
		fmt.Println("Gracefully shutdown , Doing Cleanup task.... 🧹")

		close(doneCh)

		err := app.Close(close_logger)
		if err != nil {
			log.Fatal(err)
		}

	}()

	app.StartRss(doneCh)
}
