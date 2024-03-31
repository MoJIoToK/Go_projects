package files

import (
	"article-advisor/lib/er"
	"article-advisor/storage"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
)

// Тип который реализует интерфейс. Хранит базовый путь, в какой папке хранится всё.
type Storage struct {
	basePath string
}

const defaultPerm = 0774

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	//Способ обработки ошибок
	defer func() { err = er.WrapIfErr("can't save page", err) }()

	//определение директории куда сохраняется файл
	fPath := filepath.Join(s.basePath, page.UserName)

	//Создаёт все директории в переданный путь
	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	//Создание папки и её имени
	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	//Создание файла
	file, err := os.Create(fPath)
	if err != nil {
		return err
	}

	//Искусственно сделанная конструкция, для того чтобы показать, что мы сознательно игнорируем ошибку
	defer func() { _ = file.Close() }()

	//Страница преобразовывается в формат gob и записывается в указанный файл
	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	//Способ обработки ошибок
	defer func() { err = er.WrapIfErr("can't pick random page", err) }()

	//определение директории куда был сохранен файл
	path := filepath.Join(s.basePath, userName)

	//Получение списка файлов
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return er.Wrap("can't remove file", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file %s", path)

		return er.Wrap(msg, err)
	}

	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, er.Wrap("can't check if file exists", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", path)

		return false, er.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, er.Wrap("can't decode page", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, er.Wrap("can't decode page", err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
