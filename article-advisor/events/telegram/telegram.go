package telegram

import (
	"article-advisor/clients/telegram"
	"article-advisor/events"
	"article-advisor/lib/er"
	"article-advisor/storage"
	"errors"
)

// Processor processes events collected by the Fetcher using storage.Storage.
type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

// Meta is structure for telegram bot, that stores the user number and his name
type Meta struct {
	ChatID   int
	Username string
}

// Errors
var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

// New is a constructor Processor structure.
func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

// Fetch is method for collecting events from the tg client.
func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	//Receiving updates from client.
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, er.Wrap("can't get events", err)
	}

	//Return nil, if slice of updates empty. Maybe return error about the slice is empty.
	if len(updates) == 0 {
		return nil, nil
	}

	//Slice for result.
	res := make([]events.Event, 0, len(updates))

	//iterating all the update and convert them into events.
	for _, u := range updates {
		res = append(res, event(u))
	}

	//Update offset in order to get next time updates.
	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

// Process is method that performs different actions depending on the event type.
func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return er.Wrap("can't process message", ErrUnknownEventType)
	}
}

// processMessage is method that process message in meta of event. This method calls Processor.doCMD method.
func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return er.Wrap("can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return er.Wrap("can't process message", err)
	}

	return nil
}

// meta is function that returns Meta in event.
func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, er.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

// event is function for transform update into event.
func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

// fetchText is function that transform message text of update to text of event.
func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

// fetchType is function that transform type of update to type of event.
func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}
