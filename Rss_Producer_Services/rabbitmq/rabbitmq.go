package mqueue

import (
	"context"
	"errors"
	"log"
	"time"

	mq "github.com/rabbitmq/amqp091-go"
)

type Mqclient interface {
	// Create new channel on existing connection  of MqClient
	InitializeChannel() (MqChannel, error)
	// Close the client , which include both connection and channel
	Close() error
}

type MqConnection interface {
	Close() error
}

type MqChannel interface {
	Confirm(noWait bool) error
	Consume(
		queueName string,
		config ConsumeConfig,
		args mq.Table,
	) (<-chan mq.Delivery, error)
	// Publish with no confirm
	UnsafePublish(ctx context.Context, config PublishConfig, args mq.Publishing) error
	// publish with confirm
	Publish(ctx context.Context, config PublishConfig, args mq.Publishing) error
	QueueDeclare(
		queueName string,
		config QueueDeclareConfig,
		args mq.Table,
	) (mq.Queue, error)
	QueueBind(queueName string, config QueueBindConfig, args mq.Table) error

	ExchangeDeclare(
		name string,
		config ExchangeDeclareConfig,
		args mq.Table,
	) error
	Close() error
}

type ConsumeConfig struct {
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
}

// not support table yet
type PublishConfig struct {
	ExchangeName string
	RoutingKey   string
	Mandatory    bool
	Immediate    bool
	Timeout      time.Duration

	ResendDelay time.Duration
}

type QueueDeclareConfig struct {
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
}

type QueueBindConfig struct {
	BindingKey   string
	ExchangeName string
	NoWait       bool
}

type ExchangeDeclareConfig struct {
	Type       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	Nowait     bool
}

//-----

type RabbitOpts struct {
	NoRetries      int
	ReconnectDelay time.Duration
	ReInitDelay    time.Duration
	ConfirmMode    bool
}

type rabbitMqClient struct {
	address    string
	logger     *log.Logger
	connection *mq.Connection
	channel    *mq.Channel
	done       chan bool
	opts       *RabbitOpts

	// notify when channel is ready to perform ops
	notifyReady chan bool
	// notify when channel has close  after perform ops
	notifyExit chan bool

	// notify when conn retry excess limit
	notifyConnError chan struct{}

	// notify by rabbitmq when conn close by error
	notifyConnClose chan *mq.Error

	// notify by rabbitmq when chan close by error
	notifyChanClose chan *mq.Error

	// use to ack between server and client
	notifyConfirm chan mq.Confirmation
	//
	isReady bool
}

func NewClient(address string, logger *log.Logger, opts *RabbitOpts) Mqclient {
	if opts == nil {
		opts.ReconnectDelay = 5 * time.Second
		opts.ReInitDelay = 2 * time.Second
		opts.NoRetries = 255

	}

	client := rabbitMqClient{
		opts:        opts,
		logger:      logger,
		done:        make(chan bool),
		notifyReady: make(chan bool),
		address:     address,
	}

	client.handleReconnect()

	return &client
}

// Handle Reconnect of amqp connection in case of error
func (r *rabbitMqClient) handleReconnect() {
	go func() {
		for i := 0; i < r.opts.NoRetries; i++ {
			r.isReady = false
			r.logger.Println("Conntecting RabbitMq")
			err := r.connect()
			if err != nil {
				r.logger.Printf("RabbitMq connection failed  Retrying... : %v", err)
				select {
				case <-r.done:
					return
				case <-time.After(r.opts.ReconnectDelay):
				}
				continue
			}
			i = 0
			select {
			case <-r.notifyConnError:
				continue
			case <-r.done:
				break
			}
		}
		r.logger.Fatal("RabbitMq Excess Attempt to reconnect")
	}()

}

// Create new amqp  connection
func (r *rabbitMqClient) connect() error {
	conn, err := mq.Dial(r.address)
	if err != nil {
		return err
	}
	r.logger.Println("RabbitMq connected")
	r.connection = conn
	return nil
}

// Create a channel
func (r *rabbitMqClient) InitializeChannel() (MqChannel, error) {
	for i := 0; i < r.opts.NoRetries; i++ {
		if r.connection == nil {
			select {
			case <-time.After(r.opts.ReInitDelay):
				if r.connection != nil {
					break
				}
				r.notifyConnError <- struct{}{}
				time.Sleep(r.opts.ReconnectDelay)
				continue
			case <-r.done:
				return nil, ErrShutdown
			}
			if !r.connection.IsClosed() {
				break
			}
		}
	}

	go func() {
		for i := 0; i < r.opts.NoRetries; i++ {
			r.isReady = false
			err := r.init()

			if err != nil {
				r.logger.Printf("Failed to initialize channel. Retrying... :  %v", err)
				select {
				case <-r.done:
					r.notifyConnError <- struct{}{}
				case <-r.notifyConnClose:
					r.logger.Println("Connecting closed. Reconnecting")
					r.notifyConnError <- struct{}{}
				case <-time.After(r.opts.ReInitDelay):
				}
				continue
			}
			select {
			case <-r.notifyExit:
				r.isReady = false
				return
			case <-r.done:
				return
			case <-r.notifyConnClose:
				r.logger.Println("Connecting closed. Reconnecting")
				r.notifyConnError <- struct{}{}
			case <-r.notifyChanClose:
				r.logger.Println("Channel closed , re running init....")
			}

		}
		r.logger.Fatal("RabbitMq channel Reconnect Exist Limit")
	}()

	r.logger.Println("[Dev] : Waiting for Channel")
	<-r.notifyReady

	return r.getChannel()

}

func (r *rabbitMqClient) init() error {
	ch, err := r.connection.Channel()
	if err != nil {
		return err
	}
	err = ch.Confirm(r.opts.ConfirmMode)

	if err != nil {
		return err
	}

	r.changeChannel(ch)
	r.isReady = true
	r.notifyReady <- true
	r.logger.Println("[Dev] : RabbitMq Connection ready")
	return nil

}

func (r *rabbitMqClient) changeChannel(channel *mq.Channel) {

	r.channel = channel
	r.notifyChanClose = make(chan *mq.Error, 1)
	r.notifyConfirm = make(chan mq.Confirmation, 1)
	r.notifyExit = make(chan bool, 1)

	r.channel.NotifyClose(r.notifyChanClose)
	r.channel.NotifyPublish(r.notifyConfirm)

}

// Close both existing channel and
func (r *rabbitMqClient) Close() error {
	if !r.isReady && r.connection.IsClosed() {
		return ErrAlreadyClosed
	}
	close(r.done)

	err := r.channel.Close()
	if err != nil {
		return ChannelClosingError{error: err}
	}
	err = r.connection.Close()
	if err != nil {
		return ConnectionClosingError{error: err}
	}

	r.isReady = false
	return nil

}

func (r *rabbitMqClient) getChannel() (MqChannel, error) {

	select {
	case <-r.done:
		return nil, ErrAlreadyClosed
	default:
	}
	r.logger.Println("[Dev] : Channel is ready")

	return &rabbitChannel{
		noRetries:       r.opts.NoRetries,
		done:            r.done,
		notifyConfirm:   r.notifyConfirm,
		notifyChanClose: r.notifyChanClose,
		notifyExist:     r.notifyExit,

		channel: r.channel,
	}, nil

}

type rabbitChannel struct {
	noRetries       int
	notifyExist     chan bool
	notifyConfirm   chan mq.Confirmation
	notifyChanClose chan *mq.Error
	done            chan bool
	channel         *mq.Channel
}

// TODO : Implement this
func (r *rabbitChannel) Confirm(noWait bool) error {
	return errors.New("Implement me ")
}

func (r *rabbitChannel) QueueDeclare(
	queueName string,
	config QueueDeclareConfig,
	data mq.Table,
) (mq.Queue, error) {
	return r.channel.QueueDeclare(
		queueName,
		config.Durable,
		config.AutoDelete,
		config.Exclusive,
		config.NoWait,
		data,
	)

}

func (r *rabbitChannel) Consume(
	queueName string,
	config ConsumeConfig,
	data mq.Table,
) (<-chan mq.Delivery, error) {
	return r.channel.Consume(
		queueName,
		config.Consumer,
		config.AutoAck,
		config.Exclusive,
		config.NoLocal,
		config.NoWait,
		data,
	)
}

func (r *rabbitChannel) Publish(
	ctx context.Context,
	config PublishConfig,
	data mq.Publishing,
) error {
	if config.ResendDelay == 0 {
		config.ResendDelay = time.Second * 2
	}

	for i := 0; i < r.noRetries; i++ {
		err := r.UnsafePublish(ctx, config, data)
		if err != nil {
			select {
			case <-r.done:
				return ErrShutdown
			case <-r.notifyChanClose:
				return ErrChannelNotConnected
			case <-time.After(config.ResendDelay):
			}
			continue
		}
		confirm := <-r.notifyConfirm
		if confirm.Ack {
			return nil
		}
	}
	return RetriesExcessError{where: "publish"}
}

func (r *rabbitChannel) UnsafePublish(
	ctx context.Context,
	config PublishConfig,
	data mq.Publishing,

) error {

	if config.Timeout == 0 {
		config.Timeout = time.Second * 5
	}
	c, cancel := context.WithTimeout(ctx, config.Timeout)

	defer cancel()

	return r.channel.PublishWithContext(
		c,
		config.ExchangeName,
		config.RoutingKey,
		config.Mandatory,
		config.Immediate,
		data,
	)
}

func (r *rabbitChannel) QueueBind(
	queueName string,
	config QueueBindConfig,
	args mq.Table,
) error {
	return r.channel.QueueBind(
		queueName,
		config.BindingKey,
		config.ExchangeName,
		config.NoWait,
		args,
	)
}

func (r *rabbitChannel) ExchangeDeclare(
	name string,
	config ExchangeDeclareConfig,
	args mq.Table,
) error {
	return r.channel.ExchangeDeclare(
		name,
		config.Type,
		config.Durable,
		config.AutoDelete,
		config.Internal,
		config.Nowait,
		args,
	)
}

func (r *rabbitChannel) Close() error {
	r.notifyExist <- true
	return r.channel.Close()

}
