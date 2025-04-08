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

// Storage is structure that stores the base path of the folder in which all pages is stored.
// This structure implements interface.
type Storage struct {
	basePath string
}

const defaultPerm = 0774

// New is constructor for Storage. Input data - basePath.
func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

// Save is the method for save storage.Page in a folder at the created path.
func (s Storage) Save(page *storage.Page) (err error) {
	//Processing errors.
	defer func() { err = er.WrapIfErr("can't save page", err) }()

	//Define path where the part of this path is username.
	fPath := filepath.Join(s.basePath, page.UserName)

	//Create directory with path.
	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	//Create filename in folder for define path.
	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	//Create file with filename in fPath.
	file, err := os.Create(fPath)
	if err != nil {
		return err
	}

	//An artificially made structure to show that we are deliberately ignoring the error.
	defer func() { _ = file.Close() }()

	//The page is converted to gob format and written to the specified file
	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

// PickRandom is the method that returns a random page from database.
func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	//Processing errors.
	defer func() { err = er.WrapIfErr("can't pick random page", err) }()

	//Determining the directory where the file was saved
	path := filepath.Join(s.basePath, userName)

	//Getting a list of files.
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

// Remove is the method that delete page in database.
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

// IsExists is the method that determines whether the given page exists or not.
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

// decodePage is method that decode page.
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

// fileName is method that create name of file in folder for define path using database.Hash.
// The hash is used to ensure that file names are not repeated
func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
