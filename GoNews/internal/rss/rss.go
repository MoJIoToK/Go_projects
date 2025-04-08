// Пакет rss позволяет десереализовать RSS-поток в структуру Post.

package rss

import (
	"GoNews/internal/model"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	strip "github.com/grokify/html-strip-tags-go"
)

// ErrBodyNil - ошибка указывающая на то, что тело запроса пустое.
var ErrBodyNil = errors.New("response body is nil")

// Parse - функция позволяет десериализовать RSS-поток в структуру Post. Функция возвращает слайс типа Post и ошибку.
func Parse(url string) ([]model.Post, error) {
	const operation = "GoNews.rss.Parse"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	if body == nil {
		return nil, fmt.Errorf("%s: %w", operation, ErrBodyNil)
	}

	var f model.Feed

	err = xml.Unmarshal(body, &f)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	var data []model.Post
	regex, err := regexp.Compile(`[\n]{2,}[\s]+`)
	if err != nil {
		log.Println("Failed to compile regular expression")
	}

	for _, item := range f.Channel.Items {
		var post model.Post
		post.Title = item.Title
		post.Link = item.Link
		desc := strip.StripTags(item.Desc)
		post.Content = regex.ReplaceAllString(desc, "\n")
		item.PubDate = strings.ReplaceAll(item.PubDate, ",", "")
		post.PubTime = timeConversation(item.PubDate)

		data = append(data, post)
	}

	return data, nil
}

// timeConversation - функция для конвертации времени из строки в формат time.Time.
func timeConversation(str string) time.Time {
	r, _ := utf8.DecodeLastRuneInString(str)
	if r == utf8.RuneError {
		return time.Now()
	}

	var t time.Time
	var err error
	switch {
	case unicode.IsDigit(r):
		t, err = time.Parse(time.RFC1123Z, str)
	default:
		t, err = time.Parse(time.RFC1123, str)
	}
	if err != nil {
		return time.Now()
	}
	return t
}
