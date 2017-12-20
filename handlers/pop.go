package handlers

import (
	"../initializers"
	"io"
	"net/http"
)

func Pop(w http.ResponseWriter, r *http.Request, registry *initializers.Registry) {
	b := registry.Messaging.Pop(r.URL.Query()["queue_id"][0], 10)
	io.WriteString(w, b.String())
}
