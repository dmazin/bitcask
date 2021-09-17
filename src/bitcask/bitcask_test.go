package bitcask

import (
	"testing"
)

func TestSetThenGet(t *testing.T) {
	key := "foo"
	value := "bar"

	err := set(key, value)

	// TODO why does the sample code use a string literal?
	// https://golang.org/doc/tutorial/add-a-test
	if err != nil {
		t.Fatalf(`%q`, err)
	}

	stored_value, err := get(key)

	if stored_value != value || err != nil {
		t.Fatalf(`Expected set to return %q but got %q, %v`, value, stored_value, err)
	}
}
