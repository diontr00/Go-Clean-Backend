package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"khanhanhtr/sample2/internal/random"
	"khanhanhtr/sample2/model"
	mqueue "khanhanhtr/sample2/rabbitmq"

	logStack "github.com/diontr00/logstack"
	"github.com/rabbitmq/amqp091-go"
)

var (
	ErrParsingRequest error = errors.New("Rss Request Fetching Parse Error")
)

type rssUsecase struct {
	mq     mqueue.Mqclient
	logger *logStack.Logger
}

func (r *rssUsecase) FetchAndInsert(c context.Context, request *model.RssRequest) error {

	ch, err := r.mq.InitializeChannel()

	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare("FetchAndInsert", mqueue.QueueDeclareConfig{
		Durable:    false,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
	}, nil)

	if err != nil {
		return err
	}

	req_id := random.RandomStringMust(5)

	data, err := json.Marshal(request)
	if err != nil {
		return ErrParsingRequest
	}

	err = ch.Publish(c, mqueue.PublishConfig{
		ExchangeName: "",
		RoutingKey:   "rpc_queue",
		Mandatory:    false,
		Immediate:    false,
	}, amqp091.Publishing{
		ContentType:   "application/json",
		CorrelationId: req_id,
		ReplyTo:       q.Name,
		Body:          data,
	})

	if err != nil {
		return err
	}

	go func() {
		r.handleResponse(req_id, q.Name, ch)
	}()

	return nil

}

func (r *rssUsecase) handleResponse(corrId string, queueName string, ch mqueue.MqChannel) {

	defer func() {
		err := ch.Close()
		if err != nil {
			r.logger.Error(err.Error())
		}
	}()

	msgs, err := ch.Consume(queueName, mqueue.ConsumeConfig{
		Consumer: "",
		NoWait:   false,
		AutoAck:  true,
		NoLocal:  false,
	}, nil)

	if err != nil {
		r.logger.Error(err.Error())
		// Implement method to notify unsuccessful
		return
	}

	for d := range msgs {
		if corrId == d.CorrelationId {
			// TODO: implement sms here
			var res model.RssResponse
			err := json.Unmarshal(d.Body, &res)

			if err != nil {
				r.logger.Error(err.Error())
				break
			}

			fmt.Println("DEV", res.Count)
		}
	}

}

func NewRssUsecase(
	logger *logStack.Logger,
	mq mqueue.Mqclient,
) model.RssUsecase {
	return &rssUsecase{
		mq:     mq,
		logger: logger,
	}
}
