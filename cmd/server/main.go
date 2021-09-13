package main

import (
	"log"

	"github.com/dmazin/naivedb"
	server "github.com/dmazin/naivedb/api/http"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db, err := naivedb.NewFileBackedNaiveDB("database")
	if err != nil {
		log.Fatal(err)
	}

	server.SetUpAndListen(db)
}
