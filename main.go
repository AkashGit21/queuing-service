package main

import (
	"log"
	"time"

	"github.com/AkashGit21/queuing-service/mq"
)

func main() {
	defer mq.PanicHandler()

	// Create a producer
	programming := mq.NewProducer("programming")
	motherTongue := mq.NewProducer("mother-tongue")

	// Open a queue for passing messages
	rmq, err := mq.OpenQueue(2)
	if err != nil {
		log.Fatalln(err)
	}

	// Create a consumer
	consumer1 := mq.NewConsumer("akash")
	consumer2 := mq.NewConsumer("akhil")

	// Link the consumer with desired queue
	err = rmq.AddConsumer(consumer1)
	if err != nil {
		log.Fatalln(err)
	}
	err = rmq.AddConsumer(consumer2)
	if err != nil {
		log.Fatalln(err)
	}

	// Pass different messages through different producers
	go func() {
		// Start sending data from producer 'programming'
		for ind := 0; ind < 10; ind++ {
			time.Sleep(1 * time.Second)
			err = programming.Publish(rmq, mq.NewMessage("go", "Golang is a modern programming language."))
		}
	}()

	go func() {
		// Start sending data from producer 'programming'
		for ind := 0; ind < 10; ind++ {
			time.Sleep(500 * time.Millisecond)
			err = programming.Publish(rmq, mq.NewMessage("python", "Python is the most used programming language."))
		}
	}()

	go func() {
		// Start sending data from producer 'mother-tongue'
		for ind := 0; ind < 10; ind++ {
			time.Sleep(800 * time.Millisecond)
			err = motherTongue.Publish(rmq, mq.NewMessage("hindi", "My mother tongue is Hindi."))
		}
	}()

	go func() {
		// Wait for 4 seconds to finish
		time.Sleep(4 * time.Second)

		// Stop consuming after 4 seconds
		rmq.StopConsuming()
	}()

	// Enable the listeners to consume after every second
	err = rmq.StartConsuming(1 * time.Second)
	if err != nil {
		log.Print(err)
	}
}
