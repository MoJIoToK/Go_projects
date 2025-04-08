package telegram

import (
	"article-advisor/lib/er"
	"article-advisor/storage"
	"context"
	"errors"
	"log"
	"net/url"
	"strings"
)

// Main commands for bot.
const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

// doCMD is method that calls processor methods depending on what came from the user.
func (p *Processor) doCmd(text string, chatID int, username string) error {
	//Delete space in text.
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, username)

	//command "Save". This command has no constants, it is just a link sent by the user.
	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

// savePage is method for save page.
func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = er.WrapIfErr("can't do command: save page", err) }()

	//Preparing the page to be saved.
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	//checking the availability of this page.
	isExists, err := p.storage.IsExists(context.Background(), page)
	if err != nil {
		return err
	}

	//Reply from a bot that this page is already exists.
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	//Saving a page
	if err := p.storage.Save(context.Background(), page); err != nil {
		return err
	}

	//Reply from a bot that this page is save.
	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

// sendRandom is method for send random page.
func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = er.WrapIfErr("can't do command: can't send random", err) }()

	page, err := p.storage.PickRandom(context.Background(), username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(context.Background(), page)
}

// sendHelp is method for send help message.
func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

// sendHello is method for send start message.
func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

// isAddCmd calls isURL. And returns bool data.
func isAddCmd(text string) bool {
	return isURL(text)
}

// isURL is function that determines text is a link.
func isURL(text string) bool {

	//Parsing link. Warning - the link must have a `http/https` written.
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
