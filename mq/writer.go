package mq

import (
	"log"
	"time"
)

type Producer interface {
	OpenQueue(consumersAllowed int) (*queue, error)
	Publish(q *queue, msg *message) error
}

type producer struct {
	Name string
}

// Opens a new Queue to receive and send messages
func OpenQueue(consumersAllowed int) (*queue, error) {
	return newQueue(consumersAllowed), nil
}

// Publish sends the message 'msg' over queue 'q' through producer 'p'
func (p *producer) Publish(q *queue, msg message) error {
	log.Printf("%s: publishing message", p.Name)
	// Check if queue is valid.
	if q.Data == nil {
		return ErrQueueUnavailable
	}

	if len(q.Data) > MaxNumOfMessages {
		return ErrQueueFilled
	}

	err := q.enqueueMessageAsJSON(msg)
	if err != nil {
		return err
	}

	return nil
}

func NewProducer(tag string) *producer {
	return &producer{
		Name: tag,
	}
}

type message struct {
	Created int64  `json:"created"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

func NewMessage(title string, msg string) message {
	return message{
		Created: time.Now().Unix(),
		Name:    title,
		Message: msg,
	}
}
