package events

// Fetcher collects events sent by the user. Event, for example, is different messages from the user.
// Fetcher collects events using telegram.Client.
type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

// Processor is interface with method Process.
// Processor processes events collected by the Fetcher using storage.Storage.
type Processor interface {
	Process(e Event) error
}

// Type is a type of event.
type Type int

// Event is a request from a user received using a telegram.Client.
type Event struct {
	Type Type
	Text string
	Meta interface{}
}

// List of events.
const (
	Unknown Type = iota
	Message
)
