package memory

import "github.com/Suhaan-Bhandary/pubsub"

// publisher is an in-memory generic implementation of the pubsub.Publisher interface.
// It supports publishing events of any type, specified by the Event type parameter.
type publisher[Event any] struct{}

// NewPublisher creates a new in-memory publisher for a specific Event type.
func NewPublisher[Event any]() pubsub.Publisher[Event] {
	return &publisher[Event]{}
}

func (p *publisher[Event]) Publish(event Event) error {
	panic("unimplemented")
}

func (p *publisher[Event]) Subscribe(Subscriber pubsub.Subscriber[Event]) {
	panic("unimplemented")
}

func (p *publisher[Event]) UnSubscribe(Subscriber pubsub.Subscriber[Event]) {
	panic("unimplemented")
}

func (p *publisher[Event]) Close() error {
	panic("unimplemented")
}

// subscriber is an in-memory generic implementation of the pubsub.Subscriber interface.
// It supports receiving events of any type, specified by the Event type parameter.
type subscriber[Event any] struct{}

// NewSubscriber creates a new in-memory subscriber for a specific Event type.
func NewSubscriber[Event any]() pubsub.Subscriber[Event] {
	return &subscriber[Event]{}
}

func (s *subscriber[Event]) Listen() (Event, bool) {
	panic("unimplemented")
}
