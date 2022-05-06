package mq

import "testing"

func TestNewConsumer(t *testing.T) {

}

func mockConsumer() *consumer {
	return NewConsumer("test-consumer")
}

func TestConsume(t *testing.T) {
	cns := mockConsumer()

	err := cns.Consume([]byte("some random message to be consumed"))
	if err != nil {
		t.Errorf("Unable to consume message: %v", err)
	}
}
