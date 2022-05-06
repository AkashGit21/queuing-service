package mq

import (
	"encoding/json"
	"sync"
	"time"
)

const (
	// Maximum number of messages to be handled by the queue at any time
	MaxNumOfMessages = 100
	// Default pool size of Consumers
	MaxConsumersPoolSize = 5
)

// A FIFO based Queue implementation
type Queue interface {
	StartConsuming(intervalDuration time.Duration) error
	StopConsuming() <-chan struct{}
	AddConsumer(listener *consumer) error

	// internal operation to consume messages from the queue 'q'
	consume(q *queue) error
	// internal operations for queue
	dequeue() (*message, error)
	enqueue(value []byte) error
	enqueueMessageAsJSON(msg message) error
}

type queue struct {
	sync.RWMutex
	Data chan []byte
	// list of consumers to receive data
	consumers []*consumer
	// number of messages present at the moment
	size int
	// status of Queue - open/closed
	isOpen bool
}

// newQueue generates a new Queue with required default parameters.
func newQueue(poolSize int) *queue {
	// If the poolSize is not appropriate, create default pool
	if poolSize < 1 {
		poolSize = MaxConsumersPoolSize
	}

	return &queue{
		Data:      make(chan []byte, MaxNumOfMessages),
		consumers: make([]*consumer, 0, poolSize),
		size:      0,
		isOpen:    false,
	}
}

// Enqueues the value to queue
func (q *queue) enqueue(value []byte) error {
	q.Lock()
	defer q.Unlock()

	// Check if Queue is closed.
	if !q.isOpen {
		return ErrQueueClosed
	}

	if len(q.Data) < MaxNumOfMessages && q.size < MaxNumOfMessages {
		// Push the msg to buffered channel Data of queue.
		q.Data <- value
		q.size++
	} else {
		// Return error if queue is filled.
		return ErrQueueFilled
	}

	return nil
}

// Enqueues the message into Queue
func (q *queue) enqueueMessageAsJSON(msg message) error {

	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return q.enqueue(jsonBytes)
}

// Dequeue the message from Queue
func (q *queue) dequeue() (*message, error) {
	q.Lock()
	defer q.Unlock()

	// Check if Queue is closed.
	if !q.isOpen {
		return nil, ErrQueueClosed
	}

	var msg message
	if q.size > 0 {
		value := <-q.Data
		json.Unmarshal(value, &msg)
		q.size--
	} else {
		return nil, ErrQueueEmtpy
	}
	return &msg, nil
}

// AddConsumer subscribes Consumer to the given Queue.
func (q *queue) AddConsumer(cns *consumer) error {
	q.Lock()
	defer q.Unlock()

	// Check if consumer exists
	if cns == nil {
		return ErrConsumerUnavailable
	}

	// Check if Queue is closed.
	if !q.isOpen {
		return ErrQueueClosed
	}

	if len(q.consumers) < cap(q.consumers) {
		q.consumers = append(q.consumers, cns)
	} else {
		return ErrQueueFilled
	}

	return nil
}

// Start consuming messages until the queue is empty
// Wait for the 'intervalDuration' when queue is empty before consuming again
func (q *queue) StartConsuming(intervalDuration time.Duration) error {

	// error channel to fetch errors from the goroutine
	errChan := make(chan error)

	go func(que *queue) {
		for {
			time.Sleep(intervalDuration)

			// Check if queue is closed.
			if !que.isOpen {
				errChan <- ErrQueueClosed
				return
			}

			// Check if queue is valid.
			if que.Data == nil {
				errChan <- ErrQueueUnavailable
				return
			}

			// Check if Queue is empty
			if que.size < 0 {
				// Return since the queue is empty
				return
			}

			err := consume(que)
			if err != nil {
				errChan <- err
				return
			}
		}
	}(q)

	return <-errChan
}

// StopConsuming stops consuming messages from the given Queue.
func (q *queue) StopConsuming() <-chan struct{} {
	q.Lock()
	defer q.Unlock()

	// Close the channel and resolve other attributes
	close(q.Data)
	q.size = 0
	q.isOpen = false

	return nil
}

// internal function to consume messages from Queue
func consume(q *queue) error {

	// iterate until queue is empty
	for q.size > 0 {
		msg, err := q.dequeue()
		if err != nil {
			return err
		}
		jsonBytes, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		wg := sync.WaitGroup{}
		// pass the message to every subscribed consumer
		for _, listener := range q.consumers {
			go func(c *consumer) {
				wg.Add(1)
				c.Consume(jsonBytes)
				wg.Done()
			}(listener)
		}
		// wait for every consumer to consume the message
		wg.Wait()
	}

	return nil
}
