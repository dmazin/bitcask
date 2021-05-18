package naivedb

import (
	"bytes"
	"testing"
)

func TestGetBeforeSet(t *testing.T) {
	buf := new(bytes.Buffer)
	db := NaiveDB{buf}

	key := "foo"
	stored_value, err := db.get(key)

	if err != nil {
		t.Fatalf(`%q`, err)
	}

	if stored_value != "" || err != nil {
		t.Fatalf(`Expected set to return %q but got %q, %v`, "", stored_value, err)
	}
}

func TestSetThenGet(t *testing.T) {
	buf := new(bytes.Buffer)
	db := NaiveDB{buf}

	key := "foo"
	value := "bar"

	err := db.set(key, value)

	// TODO why does the sample code use a string literal?
	// https://golang.org/doc/tutorial/add-a-test
	if err != nil {
		t.Fatalf(`%q`, err)
	}

	stored_value, err := db.get(key)

	if stored_value != value || err != nil {
		t.Fatalf(`Expected set to return %q but got %q, %v`, value, stored_value, err)
	}
}
