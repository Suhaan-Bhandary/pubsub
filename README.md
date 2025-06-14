# Go PubSub

A lightweight, type-safe, in-memory Publish/Subscribe (Pub/Sub) library for Go. This library provides a generic interface to publish and subscribe to events using custom event types.

## Features

- Generic support for any event type
- Simple `Publisher` and `Subscriber` interfaces
- In-memory implementation (currently only supports in-memory dispatch)

## Installation

```bash
go get github.com/yourusername/pubsub
````

## Usage

### Define an Event Type

```go
type MyEvent struct {
	ID   int
	Name string
}
```

### Create a Publisher

```go

var hooks = memory.Hooks[TestEvent]{
	OnPublish: func(e TestEvent) {
		fmt.Println("Published:", e)
	},
	OnSubscribe: func(sub pubsub.Subscriber[TestEvent]) {
		fmt.Println("Subscriber added", sub)
	},
	OnUnSubscribe: func(sub pubsub.Subscriber[TestEvent]) {
		fmt.Println("Subscriber removed", sub)
	},
	OnClose: func() {
		fmt.Println("Publisher closed")
	},
}

publisher := memory.NewPublisher[MyEvent](hooks)
```

### Create and Subscribe a Subscriber

```go

var opts = memory.SubscriberOptions{
	BufferSize: 10000,
}

subscriber := memory.NewSubscriber[MyEvent](opts)
publisher.Subscribe(subscriber)
```

### Publish Events

```go
event := MyEvent{ID: 1, Name: "Sample Event"}
_ = publisher.Publish(event)
```

### Receive Events

```go
for {
    e, ok := subscriber.Listen()
    if !ok {
        break
    }

    fmt.Println("Received event:", e)
}

// or

for event := range subscriber.C() {
    fmt.Println("Received event:", e)
}

```

### Clean Up

```go
// Removes a subscriber
publisher.UnSubscribe(subscriber)

// removes all subscribers and closes the publisher 
_ = publisher.Close()
```

## Event Envelope

Use `EventEnvelope` to wrap events with a name or identifier:

```go
envelope := pubsub.EventEnvelope[string, MyEvent]{
	Name: "UserCreated",
	Data: MyEvent{ID: 2, Name: "New User"},
}
```
