//Пакет model содержит модели новостей, комментариев и новостей с комментариями в базе данных.

package model

import "time"

// NewsFullDetailed - структура новости и комментария к ней.
type NewsFullDetailed struct {
	News     NewsShortDetailed
	Comments []FullComment `json:"Comments"`
}

// NewsShortDetailed - структура короткой записи новости.
type NewsShortDetailed struct {
	ID      string    `json:"ID"`
	Title   string    `json:"Title"`
	Content string    `json:"Content"`
	PubTime time.Time `json:"PubTime"`
	Link    string    `json:"Link"`
}

// Comment - структура комментария к новости.
type Comment struct {
	ID       string    `json:"ID"`
	ParentID string    `json:"ParentID"`
	NewsID   string    `json:"NewsID"`
	PubTime  time.Time `json:"PubTime"`
	Content  string    `json:"Content"`
}

// FullComment - структура комментария и его дочерних комментариев.
type FullComment struct {
	Comment
	Childs []FullComment `json:"Childs"`
}
