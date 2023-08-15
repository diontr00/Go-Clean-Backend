package mqueue

import (
	"errors"
	"fmt"
)

var (
	// Return when expected on unconnected connection
	ErrNotConnected        = errors.New("Not connected to a server")
	ErrAlreadyClosed       = errors.New("Already closed: not connected to the server")
	ErrShutdown            = errors.New("client is shuting down")
	ErrChannelNotConnected = errors.New("Channel is not connected")
)

// When closing Channel
type ChannelClosingError struct {
	error error
}

func (c ChannelClosingError) Error() string {
	return fmt.Sprintf("RabbitMQ Channel Closing Error : %v", c.error)
}

type ConnectionClosingError struct {
	error error
}

func (c ConnectionClosingError) Error() string {
	return fmt.Sprintf("RabbitMQ Connection Closing Error : %v", c.error)
}

type RetriesExcessError struct {
	where string
}

func (r RetriesExcessError) Error() string {
	return fmt.Sprintf("RabbitMQ Retries Excess Error : %s", r.where)
}
