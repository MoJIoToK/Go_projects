package storage

import (
	"article-advisor/lib/er"
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
)

// Storage is interface. Storage contains the methods for work with input data. This interface can work with files
// and databases.
type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, userName string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error)
}

var ErrNoSavedPages = errors.New("no saved page")

// Page is main data type that Storage works with. The page to which the link that was sent to the bot leads.
type Page struct {
	//URL is main parameter of this structure. URL is the link sent by the user.
	URL string
	//UserName is nickname of user that send link.
	UserName string
	//Created  time.Time
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", er.Wrap("can't calculate hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", er.Wrap("can't calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
