package event_consumer

import (
	"article-advisor/events"
	"log"
	"time"
)

//Event_consumer implements the consumer interface.

// Consumer is the structure for constantly receives and process events. Consumer performs its work using events.Fetcher
// and events.Processor. BatchSize is еру number of events processed at one time.
type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

// New is constructor of Consumer.
func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

// Start is the method that launches the Consumer.
func (c Consumer) Start() error {
	//An infinite loop that waits for events and processes them.
	for {
		// receiving an event using a fetcher.
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		//Checking for events. If count of events is 0, then we skip the iteration and wait one second.
		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		//
		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)

			continue
		}
	}
}

// handleEvents is the method for iterating and processing the event queue. Events is a slice of events.Event,
// which is received in the events.Fetcher method Fetch.
func (c *Consumer) handleEvents(events []events.Event) error {
	//sync.WaitGroup{}
	//Iterate events in the event queue.
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		//Processing the event
		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", err.Error())

			continue
		}
	}

	return nil
}
