package naivedb

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type ReaderStringWriter interface {
	io.Reader
	io.StringWriter
	io.Seeker
}

type NaiveDB struct {
	store    ReaderStringWriter
	hintStore io.ReadWriter
	offsetMap map[string]int64
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
	currentOffset, err := db.store.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	// fmt.Println(currentOffset)

	_, err = db.store.WriteString(fmt.Sprintf("%s,%s\n", key, value))
	db.offsetMap[key] = currentOffset
	return err
}

func (db *NaiveDB) get(key string) (value string, err error) {
	offset, ok := db.offsetMap[key]
	if !ok {
		panic("oh no")
	}

	db.store.Seek(offset, io.SeekStart)

	scanner := bufio.NewScanner(db.store)
	for scanner.Scan() {
		currentOffset, err := db.store.Seek(0, io.SeekCurrent)
		if err != nil {
			log.Fatalln(err)
		}
		println(currentOffset)
		line := scanner.Text()

		fmt.Printf(line)

		if strings.Contains(line, key) {
			value = strings.Split(line, ",")[1]
		}
	}

	err = scanner.Err()
	return value, err
}
