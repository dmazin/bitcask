package server

import (
	"fmt"
	"net/http"
)

type getHandler struct {
	db KVGetter
}

func (h *getHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	value, err := h.db.Get(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("error: error getting key=`%s`: %s", key, err)))

		return
	}

	fmt.Fprintf(w, "key=%s,value=%s", key, value)
}
