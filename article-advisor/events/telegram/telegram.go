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

type Meta struct {
	ChatID   int
	Username string
}

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

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	//Получение updates
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, er.Wrap("can't get events", err)
	}

	//Возвращение нулей, если список updates пуст. Можно возвращать ошибку о том, что список updates пуст
	if len(updates) == 0 {
		return nil, nil
	}

	//Подготовка слайса под результат
	res := make([]events.Event, 0, len(updates))

	//Преобразуем всех updates в event
	for _, u := range updates {
		res = append(res, event(u))
	}

	//Обновление offset для того, чтобы получить в следующий раз updates
	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return er.Wrap("can't process message", ErrUnknownEventType)
	}
}

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

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, er.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

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

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}
