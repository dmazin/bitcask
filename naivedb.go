package naivedb

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type NaiveDB struct {
	store     *os.File
	hintStore *os.File
	offsetMap map[string]int64
}

func attemptLoadOffsetMap(r io.ReadCloser, obj interface{}) (err error) {
	// decodes an arbitrary obj from r
	// todo rename me to be more general (works on more than just offsetMaps)
	// or make it a method of NaiveDB like generateOffsetMap
	dec := gob.NewDecoder(r)
	if err := dec.Decode(obj); err != nil {
		r.Close() // ignore closing error; Encode error takes precedence
		return err
	}

	log.Printf("loaded object %v", obj)

	return err
}

func attemptSaveOffsetMap(r io.WriteCloser, obj interface{}) (err error) {
	// encodes an arbitrary obj to r
	// todo rename me to be more general (works on more than just offsetMaps)
	// or make it a method of NaiveDB like generateOffsetMap
	dec := gob.NewEncoder(r)
	if err := dec.Encode(obj); err != nil {
		r.Close() // ignore closing error; Encode error takes precedence
		return err
	}

	log.Printf("saved object %v", obj)

	return err
}

func (db *NaiveDB) generateOffsetMap() (err error) {
	// when we create the hintStore, which is exactly when there
	// is nothing in the file to read.
	currentOffset, err := db.store.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	log.Printf("generating offset map. current map: %v", db.offsetMap)
	log.Printf("starting at offset %v. should be 0!", currentOffset)

	scanner := bufio.NewScanner(db.store)
	for scanner.Scan() {
		line := scanner.Text()

		split_line := strings.Split(line, ",")
		key := split_line[0]

		db.offsetMap[key] = currentOffset

		log.Printf("key=%s is at offset %v", key, currentOffset)
		currentOffset += int64(len(line))
	}

	log.Printf("generated offset map. current map: %v", db.offsetMap)

	return err
}

func (db *NaiveDB) Close() {
	db.store.Close()
	db.hintStore.Close()

	log.Printf("closed store and hintStore")
}

type NaiveDBOptions struct {
	dataPath string
}

func NewNaiveDB(options NaiveDBOptions) (_ *NaiveDB, err error) {
	storeFilepath := filepath.Join(options.dataPath, "store")

	// store is our source of truth
	// todo filename -> path
	store, err := os.OpenFile(storeFilepath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	// hintStore is a checkpoint of offsetMap so we don't have to generate it every startup
	hintStoreFilepath := filepath.Join(options.dataPath, "hintStore")
	hintStore, err := os.OpenFile(hintStoreFilepath, os.O_RDWR, 0644)

	// offsetMap tells you how many bytes from io.SeekStart you have to seek to get to the key/value pair
	offsetMap := make(map[string]int64)
	createdHintStore := false
	if err != nil {
		// Not really an error -- just means we need to create the file
		if errors.Is(err, os.ErrNotExist) {
			hintStore, err = os.Create(hintStoreFilepath)
			if err != nil {
				return nil, err
			}

			// Don't return os.ErrNotExist from main fn
			err = nil

			createdHintStore = true
		} else {
			return nil, err
		}
	}

	db := NaiveDB{store, hintStore, offsetMap}

	if createdHintStore {
		// TODO Is there some other way to assign to err, but initialize fi?
		var fi os.FileInfo
		fi, err = store.Stat()
		if err != nil {
			return nil, err
		}

		if fi.Size() > 0 {
			// TODO How can I test that this will get called only if db.store is nonempty?
			err = db.generateOffsetMap()
		}
	} else {
		attemptLoadOffsetMap(hintStore, &db.offsetMap)
	}

	return &db, err
}

func (db *NaiveDB) Set(key string, value string) (err error) {
	currentOffset, err := db.store.Seek(0, io.SeekEnd)
	if err != nil {
		db.store.Close() // ignore closing error; Seek error takes precedence
		return err
	}

	_, err = db.store.WriteString(fmt.Sprintf("%s,%s\n", key, value))
	if err != nil {
		db.store.Close() // ignore closing error; WriteString error takes precedence
		return err
	}

	log.Printf("wrote %s,%s to store at offset %v", key, value, currentOffset)

	db.offsetMap[key] = currentOffset
	attemptSaveOffsetMap(db.hintStore, db.offsetMap)

	return err
}

func (db *NaiveDB) Get(key string) (value string, err error) {
	offset := db.offsetMap[key]
	log.Printf("my offsetmap tells me that key=%s is at offset %v", key, offset)
	// fixme return an error if the key is missing

	db.store.Seek(offset, io.SeekStart)

	scanner := bufio.NewScanner(db.store)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, key) {
			value = strings.Split(line, ",")[1]
		}
	}

	err = scanner.Err()
	return value, err
}
