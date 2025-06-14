package pubsub

// Publisher represents a generic publish-subscribe message dispatcher.
// It allows publishing events and managing subscribers for a specific event type.
type Publisher[Event any] interface {
	// Publish sends an event to all subscribed listeners.
	Publish(event Event) error

	// Subscribe registers a new subscriber to receive future events.
	Subscribe(subscriber Subscriber[Event])

	// UnSubscribe removes an existing subscriber.
	UnSubscribe(subscriber Subscriber[Event])

	// Close shuts down the publisher and releases any resources.
	Close() error
}

// Subscriber represents a consumer that can listen for events of a specific type.
type Subscriber[Event any] interface {
	// Listen blocks or polls for the next available event.
	// Returns false if further events cann't be fetched.
	Listen() (Event, bool)

	// Close shuts down the Subscriber and releases any resources.
	Close()

	// Push is called by the Publisher to send an event to this subscriber.
	// User should not call this function.
	Push(event Event)

	// Acknowledge is called by the Publisher to acknowledge subscription.
	// User must not call this function.
	Acknowledge(publisher Publisher[Event])

	// AckRemoval shuts down the subscriber and releases the resources, if no publisher is listening.
	// User must not call this function.
	AckRemoval(publisher Publisher[Event]) error
}

// EventEnvelope wraps an event payload along with a named identifier.
// Name must be a string or a string-like custom type.
type EventEnvelope[Name ~string, Payload any] struct {
	Name Name
	Data Payload
}
