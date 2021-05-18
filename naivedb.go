package naivedb

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type ReaderStringWriter interface {
	io.Reader
	io.StringWriter
}

type NaiveDB struct {
	reader ReaderStringWriter
}

type FileBackedNaiveDB struct {
	db   NaiveDB
	file *os.File
}

func NewFileBackedNaiveDB(filename string) (_ *FileBackedNaiveDB, err error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	db := NaiveDB{f}
	return &FileBackedNaiveDB{db, f}, err
}

func (db *FileBackedNaiveDB) Set(key string, value string) (err error) {
	err = db.db.set(key, value)
	if err != nil {
		db.file.Close() // ignore error; Write error takes precedence
		return err
	}

	return nil
}

func (db *FileBackedNaiveDB) Get(key string) (value string, err error) {
	return db.db.get(key)
}

func (db *NaiveDB) set(key string, value string) (err error) {
	_, err = db.reader.WriteString(fmt.Sprintf("%s,%s\n", key, value))
	return err
}

func (db *NaiveDB) get(key string) (value string, err error) {
	scanner := bufio.NewScanner(db.reader)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, key) {
			value = strings.Split(line, ",")[1]
		}
	}

	err = scanner.Err()
	return value, err
}
