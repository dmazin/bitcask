package http_server

import (
	"fmt"
	"net/http"
)

type setHandler struct {
	db KVSetter
}

func (h *setHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	err := h.db.Set(key, value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("error: error setting key=`%s`,value=`%s`: %s", key, value, err)))

		return
	}

	fmt.Fprintf(w, "key=%s,value=%s", key, value)
}
