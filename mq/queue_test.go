package mq

import (
	"testing"
	"time"
)

func mockQueue() *queue {
	q, err := OpenQueue(2)
	if err != nil {
		panic(err)
	}
	return q
}

func TestOpenQueue(t *testing.T) {
	// ConsumersAllowed -> -2, 4
	q, err := OpenQueue(-2)
	if err != nil {
		t.Errorf("unable to open queue of pool size -2")
	}
	if cap(q.consumers) != MaxConsumersPoolSize {
		t.Errorf("consumers pool size should be %d", MaxConsumersPoolSize)
	}

	poolSize := 4
	q, err = OpenQueue(poolSize)
	if err != nil {
		t.Errorf("unable to open queue of pool size 4")
	}
	if cap(q.consumers) != poolSize {
		t.Errorf("consumers pool size should be %d", poolSize)
	}
}

func TestQueue(t *testing.T) {
	prod := mockProducer("test-queue")
	cns := mockConsumer()
	q := mockQueue()

	if err := q.AddConsumer(nil); err != ErrConsumerUnavailable {
		t.Errorf("nil consumer should not be added: %v", err)
	}

	if err := q.AddConsumer(cns); err != nil {
		t.Errorf("unable to add consumer: %v", err)
	}

	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			msg := NewMessage("test", "this is a test message")
			prod.Publish(q, msg)
		}
	}()

	go func() {
		time.Sleep(1 * time.Second)
		q.StopConsuming()
	}()

	q.StartConsuming(400 * time.Millisecond)
}
