package controller

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"khanhanhtr/sample2/model"
	mqueue "khanhanhtr/sample2/rabbitmq"
	"log"
	"net/http"
	"time"

	"github.com/diontr00/logstack"
	"github.com/rabbitmq/amqp091-go"
)

type rssController struct {
	logger         *logStack.Logger
	repo           model.RssRepository
	mq             mqueue.Mqclient
	retryDelay     time.Duration
	contextTimeout time.Duration
}

func (r *rssController) ListenForRssRequest(c context.Context, done chan struct{}) {
	ch, err := r.mq.InitializeChannel()

	if err != nil {
		r.logger.Fatal(err.Error())
	}

	defer func() {

		log.Fatal(ch.Close())
	}()

	q, err := ch.QueueDeclare(
		"rpc_queue",
		mqueue.QueueDeclareConfig{
			Durable:    false,
			AutoDelete: false,
			Exclusive:  false,
			NoWait:     false,
		}, nil)

	if err != nil {
		r.logger.Fatal(err.Error())
	}

	msgs, err := ch.Consume(
		q.Name,
		mqueue.ConsumeConfig{
			Consumer:  "",
			AutoAck:   false,
			Exclusive: false,
			NoLocal:   false,
			NoWait:    false,
		},
		nil,
	)

	if err != nil {
		r.logger.Fatal(err.Error())
	}

	for d := range msgs {
		d := d
		fmt.Println(string(d.Body))
		go r.processRssRequest(c, ch, d)
	}

	<-done

}

func (r *rssController) processRssRequest(
	c context.Context,
	ch mqueue.MqChannel,
	d amqp091.Delivery,
) {
	data := d.Body
	req := &model.RssRequest{}

	if err := json.Unmarshal(data, req); err != nil {
		r.handleError(ch, d, "Unable to process request")
		return
	}

	feeds, err := r.GetFeedEntries(c, req.URL)

	if err != nil {
		r.handleError(ch, d, fmt.Sprintf("Unable to fetch rss request: %v", err))
		return
	}

	count, err := r.InsertEntries(c, feeds)
	if err != nil {
		r.handleError(ch, d, fmt.Sprintf("Unable to insert rss request: %v", err))
		return
	}

	rstruc := model.RssResponse{
		InsertCount: count,
	}
	res, err := json.Marshal(&rstruc)
	if err != nil {
		r.handleError(ch, d, err.Error())
		return
	}

	if err := r.rssReply(c, ch, d, res); err != nil {
		r.logger.Fatal(err.Error())
	}

	err = d.Ack(false)
	if err != nil {
		r.logger.Error(err.Error())

	}

}

func (r *rssController) handleError(ch mqueue.MqChannel, d amqp091.Delivery, errorMessage string) {
	rstruc := model.RssResponse{
		Error: errorMessage,
	}
	res, err := json.Marshal(&rstruc)
	if err != nil {
		r.logger.Fatal(err.Error())
	}

	if err := r.rssReply(context.Background(), ch, d, res); err != nil {
		r.logger.Fatal(err.Error())
	}
}

func (r *rssController) rssReply(
	c context.Context,
	ch mqueue.MqChannel,
	d amqp091.Delivery,
	body []byte,
) error {

	err := ch.Publish(
		c,
		mqueue.PublishConfig{
			ExchangeName: "",
			RoutingKey:   d.ReplyTo,
			Timeout:      r.contextTimeout,
			Mandatory:    false,
			Immediate:    false,
			ResendDelay:  r.retryDelay,
		}, amqp091.Publishing{ContentType: "application/json", Body: body, CorrelationId: d.CorrelationId},
	)
	return err

}

func (r *rssController) InsertEntries(c context.Context, entries []*model.RSSEntry) (int, error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()
	return r.repo.InsertEntries(ctx, entries)
}

func (r *rssController) GetFeedEntries(c context.Context, url string) ([]*model.RSSEntry, error) {
	client := &http.Client{}

	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, RssRequestError{err: err}
	}
	req.Header.Add(
		"User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36(KHTML, like Gecko) Chrome/70.0.3538.110Safari/537.36",
	)

	resp, err := client.Do(req)
	if err != nil {
		return nil, RssRequestError{err: err}
	}

	if resp.StatusCode > 400 {
		return nil, RssRequestError{
			err: fmt.Errorf("Failed to fetch with status code: %v", resp.StatusCode),
		}
	}

	defer resp.Body.Close()

	byteValue, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, RssParsingError{err: err}

	}
	var feed model.RSSFeed
	err = xml.Unmarshal(byteValue, &feed)
	if err != nil {
		return nil, RssParsingError{err: err}

	}
	return feed.Entries, nil

}

func NewRssController(
	rssRepo model.RssRepository,
	logger *logStack.Logger,
	timeout time.Duration,
	mq mqueue.Mqclient,
) model.RssController {
	return &rssController{
		logger:         logger,
		repo:           rssRepo,
		mq:             mq,
		contextTimeout: timeout,
	}
}
