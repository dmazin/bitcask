package main

import (
	"log"

	"github.com/dmazin/naivedb"
	"github.com/dmazin/naivedb/api/http_server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db, err := naivedb.NewNaiveDB("database")
	if err != nil {
		log.Fatal(err)
	}

	http_server.SetUpAndListen(db)

	defer db.Close()
}
