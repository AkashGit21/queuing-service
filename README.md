# Queuing Service

## Components
  - **Queue**
  - **Writer** (Producer)
  - **Listener** (Consumer)

## Features
  - Enqueue a JSON document into the queue. **enqueue()**
  - Dequeue a JSON document from the queue. **dequeue()**
  - Multiple producers (writers) can concurrently write in the queue. **Publish()**
  - Multiple consumers (listeners) can concurrently read(listen) from the queue. **StartConsuming()** & **StopConsuming()**
  - A configurable waiting time interval for the listeners when the queue is empty, before trying again to read from the queue.
  - Message generation using **NewMessage()** with the below format: 
  ```sh
  {
    "created" : 1651818414,
    "name"    : "go",
    "message" : "Golang is a modern programming language."
  }
  ```
  - Unit and integration tests

## Usage
  - Initialize all the components, i.e: Producer, Consumer or a Pool of Producers, Consumers. This can be done using `NewProducer()` and `NewConsumer()` functions in the `mq` package.
  - Open a queue with the desired pool size of consumers using `OpenQueue()` in mq package.
  - Subscribe/Add the desired set of consumers to the previously opened queue using `AddConumer()`.
  - Once the queue is open and subscribed to, start transmitting messages using `Publish()` and those messages can be consumed periodically using `StartConsuming()`.
  - At any point of time, consuming messages from the queue can be stopped via `StopConsuming()`.

  **Note**: A use-case has been shown in the main.go following the above steps.

## References/Sources

- Medium Article: [The Stuff That Every Developer Should Know About Message Queues](https://medium.com/event-driven-utopia/the-stuff-that-every-developer-should-know-about-message-queues-a9452ac9c9d)
- Github repository: [Message queue system written in Go](https://github.com/adjust/rmq)
- Terraform: [Source code](https://github.com/hashicorp/terraform/tree/v1.1.9)
