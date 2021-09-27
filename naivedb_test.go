package naivedb

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestGetBeforeSet(t *testing.T) {
	tempDirName := t.TempDir()

	NaiveDBOptions := NaiveDBOptions{
		dataPath: tempDirName,
	}

	db, err := NewNaiveDB(NaiveDBOptions)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(db.Close)

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
	tempDirName := t.TempDir()

	NaiveDBOptions := NaiveDBOptions{
		dataPath: tempDirName,
	}

	db, err := NewNaiveDB(NaiveDBOptions)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(db.Close)

	test_data := map[string]string{
		"foo":  "bar",
		"fizz": "bazz",
		"baz":  "bat",
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

func TestGenerateOffsetMapFromDatabase(t *testing.T) {
	// First, copy the source database file to a temporary directory At first I
	// wanted to refactor everything so that NaiveDB didn't depend on os.File,
	// but [Prometheus' tests follow this same
	// pattern](https://github.com/prometheus/prometheus/blob/main/tsdb/repair_test.go#L73)
	// and I think it works out for them. Maybe file-backed dbs are a special
	// case where there is no way/reason to decouple from files
	tempDirName := t.TempDir()
	storeFilepath := filepath.Join(tempDirName, storeFilename)

	testDataStoreFilepath := filepath.Join("test_data", storeFilename)
	CopyFile(testDataStoreFilepath, storeFilepath)

	// Tests that NaiveDB.offsetMap is generated correctly from an existing database file
	db, err := NewNaiveDB(NaiveDBOptions{dataPath: tempDirName})

	// NewNaiveDB will call generateOffsetMap
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(db.Close)

	expected := map[string]int64{
		"foo":  0,
		"fizz": 7,
		"baz":  16,
	}

	if len(db.offsetMap) != len(expected) {
		t.Fatalf(`Expected offsetMap to have %v keys but got %v`, len(expected), len(db.offsetMap))
	}

	for k, v := range db.offsetMap {
		if v != expected[k] {
			t.Fatalf(`Expected offsetMap[%q] to be %q but got %q`, k, expected[k], v)
		}
	}
}

// from https://git.mills.io/prologic/bitcask/src/branch/master/bitcask_test.go#L2112
type benchmarkTestCase struct {
	name string
	size int
}

func BenchmarkGet(b *testing.B) {
	SuppressLogs(b)

	tempDirName := b.TempDir()

	NaiveDBOptions := NaiveDBOptions{
		dataPath: tempDirName,
	}

	tests := []benchmarkTestCase{
		{"128B", 128},
		{"256B", 256},
		{"512B", 512},
		{"1K", 1024},
		{"2K", 2048},
		{"4K", 4096},
		{"8K", 8192},
		{"16K", 16384},
		{"32K", 32768},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.SetBytes(int64(tt.size))

			// key := []byte("foo")
			// value := []byte(strings.Repeat(" ", tt.size))
			key := "foo"
			value := strings.Repeat(" ", tt.size)

			db, err := NewNaiveDB(NaiveDBOptions)
			if err != nil {
				b.Fatal(err)
			}

			err = db.Set(key, value)
			if err != nil {
				b.Fatal(err)
			}

			b.Cleanup(db.Close)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				val, err := db.Get(key)
				if err != nil {
					b.Fatal(err)
				}
				if val != value {
					b.Errorf("unexpected value")
				}
			}
			b.StopTimer()
		})
	}
}
