package mq

import "testing"

func init() {
}

func TestNewProducer(t *testing.T) {
	producer := NewProducer("testing-prod")
	if producer == nil {
		t.Errorf("No producer created")
	}
}

func mockProducer(tag string) *producer {
	return NewProducer(tag)
}

func mockMessage() message {
	title := "Golang"
	msg := "Go is a statically typed, compiled programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson."

	return NewMessage(title, msg)
}

func TestPublish(t *testing.T) {
	defer PanicHandler()

	type Case struct {
		name        string
		msg         message
		publisher   *producer
		que         *queue
		expectedErr error
	}

	testcases := []Case{
		{
			name:        "queue_unavailable",
			msg:         mockMessage(),
			publisher:   mockProducer("publish-test"),
			que:         nil,
			expectedErr: ErrQueueUnavailable,
		},
		{
			name:        "queue_closed",
			msg:         mockMessage(),
			publisher:   mockProducer("publish-test"),
			que:         newQueue(2),
			expectedErr: ErrQueueClosed,
		},
		{
			msg:         mockMessage(),
			publisher:   mockProducer("publish-test"),
			que:         mockQueue(),
			expectedErr: nil,
		},
	}

	for _, tcase := range testcases {
		t.Run(tcase.name, func(t *testing.T) {

			err := tcase.publisher.Publish(tcase.que, tcase.msg)
			if err != tcase.expectedErr {
				t.Errorf("Publish failed. \nwant: %v \ngot: %v", tcase.expectedErr, err)
			}

		})
	}
}
