package memdb

import "GoNews/internal/model"

type Storage struct {
	news []model.Post
	id   int
}

// New - конструктор для эмулятора подключения к ДБ
func New() *Storage {
	return &Storage{
		news: posts,
	}
}

// GetPosts - метод возвращает слайс с публикациями
func (s *Storage) GetPosts(n int) ([]model.Post, error) {
	var posts []model.Post
	//
	//for _, post := range s.news {
	//	posts = append(posts, post)
	//}
	for i := 0; i <= n; i++ {
		posts = append(posts, s.news[i])
	}

	return posts, nil
}

// AddPost - метод добавляет публикацию в память
func (s *Storage) AddPost(post model.Post) (int, error) {
	post.ID = len(s.news) + 1
	s.news = append(s.news, post)
	return len(s.news), nil
}

// Close - эмуляция закрытия БД
func (s *Storage) Close() error {
	return nil
}

// Len - метод возвращает количество публикаций
func (s *Storage) Len() int {
	return len(s.news)
}

var posts = []model.Post{
	{
		ID:      1,
		Title:   "Effective Go",
		Content: "Go is a new language.",
	},
	{
		ID:      2,
		Title:   "The Go Memory Model",
		Content: "The Go memory model specifies the conditions under which",
	},
}
