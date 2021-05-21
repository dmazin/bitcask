package naivedb

import (
	"bufio"
	"encoding/gob"
	"errors"
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
	store     ReaderStringWriter
	hintStore io.ReadWriter
	offsetMap map[string]int64
}

type FileBackedNaiveDB struct {
	db        NaiveDB
	store     io.Closer
	hintStore io.Closer
}

func attemptLoadOffsetMap(r io.Reader, obj interface{}) {
	// todo rename me to be more general
	dec := gob.NewDecoder(r)
	if err := dec.Decode(obj); err != nil {
		log.Fatalln(err)
	}

	log.Printf("loaded object %v", obj)
}

func attemptSaveOffsetMap(r io.Writer, obj interface{}) {
	// todo rename me to be more general
	dec := gob.NewEncoder(r)
	if err := dec.Encode(obj); err != nil {
		log.Fatalln(err)
	}

	log.Printf("saved object %v", obj)
}

func NewFileBackedNaiveDB(filename string) (_ *FileBackedNaiveDB, err error) {
	store, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	hintStoreFilename := fmt.Sprintf("%s.hint", filename)
	hintStore, err := os.OpenFile(hintStoreFilename, os.O_RDWR, 0644)
	offsetMap := make(map[string]int64)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			hintStore, err = os.Create(hintStoreFilename)
			if err != nil {
				// fixme just return these errs instead
				log.Fatalln(err)
			}

			// Don't return os.ErrNotExist from main fn
			err = nil
		} else {
			log.Fatalln(err)
		}
	} else {
		attemptLoadOffsetMap(hintStore, &offsetMap)
	}

	db := NaiveDB{store, hintStore, offsetMap}
	return &FileBackedNaiveDB{db, store, hintStore}, err
}

func (db *FileBackedNaiveDB) Set(key string, value string) (err error) {
	err = db.db.set(key, value)
	if err != nil {
		db.store.Close() // ignore error; Write error takes precedence
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

	_, err = db.store.WriteString(fmt.Sprintf("%s,%s\n", key, value))
	db.offsetMap[key] = currentOffset

	attemptSaveOffsetMap(db.hintStore, db.offsetMap)

	return err
}

func (db *NaiveDB) get(key string) (value string, err error) {
	offset := db.offsetMap[key]
	// fixme return an error if the key is missing

	db.store.Seek(offset, io.SeekStart)

	scanner := bufio.NewScanner(db.store)
	for scanner.Scan() {
		_, err := db.store.Seek(0, io.SeekCurrent)
		if err != nil {
			log.Fatalln(err)
		}

		line := scanner.Text()

		if strings.Contains(line, key) {
			value = strings.Split(line, ",")[1]
		}
	}

	err = scanner.Err()
	return value, err
}
