package pubsub

// Publisher represents a generic publish-subscribe message dispatcher.
// It allows publishing events and managing subscribers for a specific event type.
type Publisher[Event any] interface {
	// Publish sends an event to all subscribed listeners.
	Publish(event Event) error

	// Subscribe registers a new subscriber to receive future events.
	Subscribe(Subscriber Subscriber[Event])

	// UnSubscribe removes an existing subscriber.
	UnSubscribe(Subscriber Subscriber[Event])

	// Close shuts down the publisher and releases any resources.
	Close() error
}

// Subscriber represents a consumer that can listen for events of a specific type.
type Subscriber[Event any] interface {
	// Listen blocks or polls for the next available event.
	// Returns false if further events cann't be fetched.
	Listen() (Event, bool)
}

// EventEnvelope wraps an event payload along with a named identifier.
// Name must be a string or a string-like custom type.
type EventEnvelope[Name ~string, Payload any] struct {
	Name Name
	Data Payload
}
