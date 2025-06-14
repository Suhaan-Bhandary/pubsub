package memory_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/Suhaan-Bhandary/pubsub"
	"github.com/Suhaan-Bhandary/pubsub/memory"
)

// TestEvent contains data of type string and name of type whose underlying type is string
type TestEvent pubsub.EventEnvelope[string, int]

var opts = memory.SubscriberOptions{
	BufferSize: 10000,
}

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

// TestBasicExample tests the basic scenario with one publisher and one subscriber
func TestBasicExample(t *testing.T) {
	publisher := memory.NewPublisher(hooks)
	subscriber := memory.NewSubscriber[TestEvent](opts)

	publisher.Subscribe(subscriber)

	// Go routine to publish events
	go func() {
		defer publisher.Close()
		for i := range 3 {
			publisher.Publish(TestEvent{
				Name: "TEST",
				Data: i,
			})
		}
	}()

	expected := 0
	for event := range subscriber.C() {
		if expected != event.Data {
			t.Fatalf("error listening to socket, expected: %d, got: %d", expected, event.Data)
		}

		expected++
	}
}

// TestMultipleSubscriber check working of one publisher with multiple subscribers
func TestMultipleSubscriber(t *testing.T) {
	publisher := memory.NewPublisher(hooks)

	subscriber1 := memory.NewSubscriber[TestEvent](opts)
	subscriber2 := memory.NewSubscriber[TestEvent](opts)

	publisher.Subscribe(subscriber1)
	publisher.Subscribe(subscriber2)

	// Go routine to publish events
	go func() {
		defer publisher.Close()
		for i := range 3 {
			publisher.Publish(TestEvent{
				Name: "TEST",
				Data: i,
			})
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		expected := 0
		for event := range subscriber1.C() {
			if expected != event.Data {
				t.Errorf("error listening to socket, expected: %d, got: %d", expected, event.Data)
			}
			expected++
		}
	}()

	go func() {
		defer wg.Done()

		expected := 0
		for event := range subscriber2.C() {
			if expected != event.Data {
				t.Errorf("error listening to socket, expected: %d, got: %d", expected, event.Data)
			}
			expected++
		}
	}()

	wg.Wait()
}

// TestMultiplePublisher check working of multiple publisher with single subscriber
func TestMultiplePublisher(t *testing.T) {
	publisher1 := memory.NewPublisher(hooks)
	publisher2 := memory.NewPublisher(hooks)

	subscriber := memory.NewSubscriber[TestEvent](opts)
	publisher1.Subscribe(subscriber)
	publisher2.Subscribe(subscriber)

	// Go routine to publish events
	go func() {
		defer publisher1.Close()

		for i := range 3 {
			publisher1.Publish(TestEvent{
				Name: "TEST-1",
				Data: i,
			})
		}
	}()

	go func() {
		defer publisher2.Close()

		for i := range 3 {
			publisher2.Publish(TestEvent{
				Name: "TEST-2",
				Data: i,
			})
		}
	}()

	readCount := 0
	for range subscriber.C() {
		readCount++
	}

	expectedReadCount := 6
	if readCount != expectedReadCount {
		t.Errorf("error listening to socket, expected: %d, got: %d", expectedReadCount, readCount)
	}
}

// TestMultiplePublishserSubscriber check working of multiple publisher with multiple subscribers
func TestMultiplePublishserSubscriber(t *testing.T) {
	publisher1 := memory.NewPublisher(hooks)
	publisher2 := memory.NewPublisher(hooks)

	subscriber1 := memory.NewSubscriber[TestEvent](opts)
	subscriber2 := memory.NewSubscriber[TestEvent](opts)

	publisher1.Subscribe(subscriber1)
	publisher1.Subscribe(subscriber2)

	publisher2.Subscribe(subscriber1)
	publisher2.Subscribe(subscriber2)

	// Go routine to publish events
	go func() {
		defer publisher1.Close()

		for i := range 3 {
			publisher1.Publish(TestEvent{
				Name: "TEST-1",
				Data: i,
			})
		}
	}()

	go func() {
		defer publisher2.Close()

		for i := range 3 {
			publisher2.Publish(TestEvent{
				Name: "TEST-2",
				Data: i,
			})
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		readCount := 0
		for range subscriber1.C() {
			readCount++
		}

		expectedReadCount := 6
		if readCount != expectedReadCount {
			t.Errorf("error listening to socket, expected: %d, got: %d", expectedReadCount, readCount)
		}
	}()

	go func() {
		defer wg.Done()

		readCount := 0
		for range subscriber2.C() {
			readCount++
		}

		expectedReadCount := 6
		if readCount != expectedReadCount {
			t.Errorf("error listening to socket, expected: %d, got: %d", expectedReadCount, readCount)
		}
	}()

	wg.Wait()
}
