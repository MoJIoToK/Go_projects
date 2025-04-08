//Пакет model содержит структуру для работы с запросами.

package model

// Request - структура запроса на цензурирование комментария.
type Request struct {
	Content string `json:"content"`
}
