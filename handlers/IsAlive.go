package handlers

import (
	"persistentQueue/initializers"
	"io"
	"net/http"
)

func IsAlive(w http.ResponseWriter, r *http.Request, registry *initializers.Registry) {
	io.WriteString(w, "OK")
}
