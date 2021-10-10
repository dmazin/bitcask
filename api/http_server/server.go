package http_server

import (
	"log"
	"net/http"
)

// func (db *Bitcask) Get(key string) (value string, err error) {
// 	func (db *Bitcask) Set(key string, value string) (err error) {
type KVSetter interface {
	Set(key string, value string) error
}

type KVGetter interface {
	Get(key string) (string, error)
}

type KVSetterGetter interface {
	KVSetter
	KVGetter
}

func SetUpAndListen(db KVSetterGetter) {
	http.Handle("/set", &setHandler{db})
	http.Handle("/get", &getHandler{db})

	log.Println("listening!")
	http.ListenAndServe(":8080", nil)
}
