package memory_test

import (
	"testing"

	"github.com/Suhaan-Bhandary/pubsub"
)

// TestEvent contains data of type string and name of type whose underlying type is string
type TestEvent pubsub.EventEnvelope[string, string]

// TestBasicExample tests the basic scenario with one publisher and one subscriber
func TestBasicExample(t *testing.T) {
	t.Fatal("unimplemented")
}
