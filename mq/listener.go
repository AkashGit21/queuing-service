package mq

import (
	"fmt"
)

const (
	ConsumeMessageFormat = `
	Consumer: %s
	Message: %+v
	`
)

type Consumer interface {
	Consume() error
}

type consumer struct {
	Name    string
	Message []byte
}

// NewConsumer creates a new consumer with the given name
func NewConsumer(tag string) *consumer {
	return &consumer{
		Name: tag,
	}
}

// consumes the received message and display it in desired format
func (c *consumer) Consume(value []byte) error {

	c.Message = value
	fmt.Printf(ConsumeMessageFormat, c.Name, string(c.Message))

	return nil
}
