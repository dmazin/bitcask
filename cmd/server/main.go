package main

import (
	"log"

	"github.com/dmazin/bitcask"
	"github.com/dmazin/bitcask/api/http_server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db, err := bitcask.NewBitcask("database")
	if err != nil {
		log.Fatal(err)
	}

	http_server.SetUpAndListen(db)

	defer db.Close()
}
