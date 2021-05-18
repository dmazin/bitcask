package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dmazin/naivedb"
)

func main() {
	flag.Parse()

	var key, value string
	if flag.NArg() == 1 {
		key = flag.Arg(0)
	} else {
		key = flag.Arg(0)
		value = flag.Arg(1)
	}

	db, err := naivedb.NewFileBackedNaiveDB("database")
	if err != nil {
		log.Fatal(err)
	}

	if len(value) < 1 {
		value, err = db.Get(key)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		if err = db.Set(key, value); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("%s=%s", key, value)
}
