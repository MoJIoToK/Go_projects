package rss

import (
	"testing"
)

func TestParse(t *testing.T) {
	feed, err := Parse("https://habr.com/ru/rss/best/daily/?fl=ru")
	if err != nil {
		t.Fatal(err)
	}

	if len(feed) == 0 {
		t.Fatal("empty feed")
	}
}
