package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"url-shortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

// Конструктор для драйвера БД
func New(storagePath string) (*Storage, error) {
	//Для простоты нахождения ошибки. Т.е. ошибку надо будет искать в пакете database.sqlite в функции New
	const op = "database.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	//TODO: миграции БД

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	alias TEXT NOT NULL UNIQUE,
    	url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("#{op} : #{err}")
	}

	//Запуск запроса
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	return &Storage{db: db}, nil
}

// Сохранение ссылки. Возвращает индекс созданной записи или ошибку
func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	const op = "database.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s : %w", op, err)
	}

	//Запуск запроса
	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		//TODO: SQLite - Constraints
		//Проверка на наличие ссылки в бд
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("#{op} : #{database.ErrURLExists}")
		}
		return 0, fmt.Errorf("%s : %w", op, err)
	}

	//Получение id ссылки из БД
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s : failed to get last insert id: %w", op, err)
	}

	return id, nil
}

// Получение ссылки по alias
func (s *Storage) GetURL(alias string) (string, error) {
	const op = "database.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias=?")
	if err != nil {
		return "", fmt.Errorf("%s : %w", op, err)
	}

	//Подготовка переменной в которую будет записан результат
	var resultURL string

	//Поиск ссылки по alias в БД
	err = stmt.QueryRow(alias).Scan(&resultURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}

	if err != nil {
		return "", fmt.Errorf("%s : execute statement: %w", op, err)
	}

	return resultURL, nil
}

// DeleteURL - удаление ссылки по alias.
func (s *Storage) DeleteURL(alias string) error {
	const op = "database.sqlite.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias=?")
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}

	res, err := stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s : execute statement: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}

	if rowsAffected == 0 {
		return storage.ErrURLNotFound
	}

	return nil
}
