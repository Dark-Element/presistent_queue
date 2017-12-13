package handlers

import (
	"net/http"
	"presistentQueue/initializers"
	"io"
)

func IsAlive(w http.ResponseWriter, r *http.Request, registry *initializers.Registry) {
	io.WriteString(w, "OK")
}
