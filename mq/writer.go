package mq

import (
	"log"
	"time"
)

type Producer interface {
	Publish(q *queue, msg *message) error
}

type producer struct {
	Name string
}

// Opens a new Queue to receive and send messages
func OpenQueue(consumersAllowed int) (*queue, error) {
	// generate new queue
	que := newQueue(consumersAllowed)
	// open it for receiving and sending messages
	que.isOpen = true

	return que, nil
}

// Publish sends the message 'msg' over queue 'q' through producer 'p'
func (p *producer) Publish(q *queue, msg message) error {
	log.Printf("%s: publishing message", p.Name)

	// Check if queue is available.
	if q == nil || q.Data == nil {
		return ErrQueueUnavailable
	}

	// enqueue the message as JSON
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
