package naivedb

import (
	"os"
	"testing"
)

func TestGetBeforeSet(t *testing.T) {
	f, err := os.CreateTemp("", "naivedb_test")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(f.Name())

	db, err := NewFileBackedNaiveDB(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	key := "foo"
	stored_value, err := db.Get(key)

	if err != nil {
		t.Fatalf(`%q`, err)
	}

	if stored_value != "" || err != nil {
		t.Fatalf(`Expected set to return %q but got %q, %v`, "", stored_value, err)
	}
}

func TestSetThenGet(t *testing.T) {
	f, err := os.CreateTemp("", "naivedb_test")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(f.Name())

	db, err := NewFileBackedNaiveDB(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	key := "foo"
	value := "bar"

	err = db.Set(key, value)

	// TODO why does the sample code use a string literal?
	// https://golang.org/doc/tutorial/add-a-test
	if err != nil {
		t.Fatalf(`%q`, err)
	}

	stored_value, err := db.Get(key)

	if stored_value != value || err != nil {
		t.Fatalf(`Expected set to return %q but got %q, %v`, value, stored_value, err)
	}
}
