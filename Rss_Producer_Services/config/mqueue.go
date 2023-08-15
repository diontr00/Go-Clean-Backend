package config

import (
	"fmt"
	mqueue "khanhanhtr/sample2/rabbitmq"
	"log"
	"os"
)

func newMessageQueue(env *Env) *mqueue.Mqclient {
	// Todo: replace over log stack
	logger := log.New(os.Stdout, "[RSS-Services]", log.LstdFlags)
	client := mqueue.NewClient(env.Mq.MqUri, logger, &mqueue.RabbitOpts{
		NoRetries:      env.Mq.NoRetries,
		ConfirmMode:    env.Mq.ConfirmMode,
		ReInitDelay:    env.Mq.ReInitDelay,
		ReconnectDelay: env.Mq.ReConnectDelay,
	})
	fmt.Println(env.Mq.MqUri)
	return &client

}
