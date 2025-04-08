//Пакет model содержит структуры для работы с RSS-лентами и моделью публикаций в базе данных.

package model

import (
	"encoding/xml"
	"time"
)

// Post - модель публикации в БД, получаемая из RSS ленты.
type Post struct {
	ID      string    `json:"id" bson:"_id"`
	Title   string    `json:"title" bson:"title"`
	Content string    `json:"content" bson:"content"`
	PubTime time.Time `json:"pubTime" bson:"pubTime"`
	Link    string    `json:"link" bson:"link"`
}

// Feed - структура для RSS канала.
type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

// Channel - структура для RSS-потока.
type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

// Item - структура для записи одной публикации.
type Item struct {
	Title   string `xml:"title"`
	Desc    string `xml:"description"`
	PubDate string `xml:"pubDate"`
	Link    string `xml:"link"`
}
