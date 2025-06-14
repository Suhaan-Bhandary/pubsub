package memory

import (
	"fmt"
	"sync"

	"github.com/Suhaan-Bhandary/pubsub"
)

type Empty struct{}

// publisher is an in-memory generic implementation of the pubsub.Publisher interface.
// It supports publishing events of any type, specified by the Event type parameter.
type publisher[Event any] struct {
	mu          sync.RWMutex
	subscribers map[internalSubscriber[Event]]Empty
	hooks       Hooks[Event]
}

type Hooks[Event any] struct {
	OnPublish     func(event Event)
	OnSubscribe   func(subscriber pubsub.Subscriber[Event])
	OnUnSubscribe func(subscriber pubsub.Subscriber[Event])
	OnClose       func()
}

// NewPublisher creates a new in-memory publisher for a specific Event type.
func NewPublisher[Event any](hooks Hooks[Event]) pubsub.Publisher[Event] {
	return &publisher[Event]{
		subscribers: make(map[internalSubscriber[Event]]Empty),
		hooks:       hooks,
	}
}

func (p *publisher[Event]) Publish(event Event) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.hooks.OnPublish != nil {
		p.hooks.OnPublish(event)
	}

	for subscriber := range p.subscribers {
		subscriber.push(event)
	}

	return nil
}

func (p *publisher[Event]) Subscribe(subscriber pubsub.Subscriber[Event]) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.hooks.OnSubscribe != nil {
		p.hooks.OnSubscribe(subscriber)
	}

	internalSubscriber, ok := subscriber.(internalSubscriber[Event])
	if !ok {
		return fmt.Errorf("invalid subscriber")
	}

	p.subscribers[internalSubscriber] = Empty{}
	internalSubscriber.acknowledge(p)

	return nil
}

func (p *publisher[Event]) UnSubscribe(subscriber pubsub.Subscriber[Event]) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.hooks.OnUnSubscribe != nil {
		p.hooks.OnUnSubscribe(subscriber)
	}

	internalSubscriber, ok := subscriber.(internalSubscriber[Event])
	if !ok {
		return fmt.Errorf("invalid subscriber")
	}

	delete(p.subscribers, internalSubscriber)
	internalSubscriber.ackRemoval(p)

	return nil
}

func (p *publisher[Event]) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.hooks.OnClose != nil {
		p.hooks.OnClose()
	}

	for subscriber := range p.subscribers {
		subscriber.ackRemoval(p)
	}

	p.subscribers = make(map[internalSubscriber[Event]]Empty)
	return nil
}

type internalSubscriber[Event any] interface {
	pubsub.Subscriber[Event]

	// push is called by the Publisher to send an event to this subscriber.
	push(event Event)

	// acknowledge is called by the Publisher to acknowledge subscription.
	acknowledge(pub pubsub.Publisher[Event])

	// ackRemoval shuts down the subscriber and releases the resources, if no publisher is listening.
	ackRemoval(pub pubsub.Publisher[Event]) error
}

// subscriber is an in-memory generic implementation of the pubsub.Subscriber interface.
// It supports receiving events of any type, specified by the Event type parameter.
type subscriber[Event any] struct {
	eventCh chan Event

	mu         sync.RWMutex
	closed     bool
	publishers map[pubsub.Publisher[Event]]Empty
}

type SubscriberOptions struct {
	BufferSize int
}

// NewSubscriber creates a new in-memory subscriber for a specific Event type.
func NewSubscriber[Event any](opts SubscriberOptions) pubsub.Subscriber[Event] {
	bufferSize := opts.BufferSize
	if bufferSize <= 0 {
		bufferSize = 10000
	}

	return &subscriber[Event]{
		eventCh:    make(chan Event, bufferSize),
		publishers: make(map[pubsub.Publisher[Event]]Empty),
	}
}

func (s *subscriber[Event]) Listen() (Event, bool) {
	event, ok := <-s.eventCh
	return event, ok
}

func (s *subscriber[Event]) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.closed {
		close(s.eventCh)
		s.closed = true

		for publisher := range s.publishers {
			publisher.UnSubscribe(s)
		}
	}
}

func (s *subscriber[Event]) ackRemoval(publisher pubsub.Publisher[Event]) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.publishers, publisher)

	// Closing the event channel if all publisher are removed
	if len(s.publishers) == 0 && !s.closed {
		close(s.eventCh)
		s.closed = true
	}

	return nil
}

func (s *subscriber[Event]) push(event Event) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return
	}

	// Drops the event if buffer overflows
	select {
	case s.eventCh <- event:
	default:
	}
}

func (s *subscriber[Event]) acknowledge(publisher pubsub.Publisher[Event]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.publishers[publisher] = Empty{}
}
