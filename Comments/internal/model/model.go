//Пакет model содержит структуры для модели комментариев в базе данных.

package model

import "time"

// Comment - модель комментария в БД.
type Comment struct {
	ID       string    `json:"ID" bson:"_id"`
	ParentID string    `json:"ParentID" bson:"ParentID"`
	NewsID   string    `json:"NewsID" bson:"NewsID"`
	Content  string    `json:"Content" bson:"Content"`
	PubTime  time.Time `json:"PubTime" bson:"PubTime"`
	Childs   []Comment `json:"Childs" bson:"Childs"`
}
