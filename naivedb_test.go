package naivedb

import (
	"os"
	"testing"
)

// func TestGetBeforeSet(t *testing.T) {
// 	f, err := os.CreateTemp("", "naivedb_test")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	defer os.Remove(f.Name())

// 	db, err := NewNaiveDB(f.Name())
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	defer db.Close()

// 	key := "foo"
// 	stored_value, err := db.Get(key)

// 	if err != nil {
// 		t.Fatalf(`%q`, err)
// 	}

// 	if stored_value != "" || err != nil {
// 		t.Fatalf(`Expected set to return %q but got %q, %v`, "", stored_value, err)
// 	}
// }

func TestSetThenGet(t *testing.T) {
	tempDirName := os.TempDir()

	// TODO defer os.Remove(f.Name())

	NaiveDBOptions := NaiveDBOptions{
		dataPath : tempDirName,
	}

	db, err := NewNaiveDB(NaiveDBOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	test_data := map[string]string{
		"foo": "bar",
		"fizz": "bazz",
		"baz": "bat",
	}
	
	for k, v := range test_data {
		err = db.Set(k, v)

		if err != nil {
			// TODO why does the sample code use a string literal?
			// https://golang.org/doc/tutorial/add-a-test
			t.Fatalf(`%q`, err)
		}

		stored_value, err := db.Get(k)
		if err != nil {
			// TODO why does the sample code use a string literal?
			// https://golang.org/doc/tutorial/add-a-test
			t.Fatalf(`%q`, err)
		}

		if stored_value != v {
			t.Fatalf(`Expected set to return %q but got %q, %v`, v, stored_value, err)
		}
	}
}

// func TestGenerateOffsetMapFromDatabase(t *testing.T) {
// 	// Tests that NaiveDB.offsetMap is generated correctly from an existing database file
// 	db, err := NewNaiveDB("test_data/database")

// 	// NewNaiveDB will call generateOffsetMap
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	defer db.Close()
// 	defer os.Remove("test_data/database.hint")

// 	expected := map[string]int64{
// 		"foo": 0,
// 		"fizz": 7,
// 		"baz": 16,
// 	}

// 	if len(db.offsetMap) != len(expected) {
// 		t.Fatalf(`Expected offsetMap to have %v keys but got %v`, len(expected), len(db.offsetMap))
// 	}

// 	for k, v := range db.offsetMap {
// 		if v != expected[k] {
// 			t.Fatalf(`Expected offsetMap[%q] to be %q but got %q`, k, expected[k], v)
// 		}
// 	}
// }
