package main

import (
	"flag"
	"log"

	"github.com/dmazin/naivedb"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Parse()

	var key, value string
	if flag.NArg() == 1 {
		key = flag.Arg(0)
	} else {
		key = flag.Arg(0)
		value = flag.Arg(1)
	}

	db, err := naivedb.NewNaiveDB("database")
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

	log.Printf("and your key/value is... %s=%s", key, value)
}
