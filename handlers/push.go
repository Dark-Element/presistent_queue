package handlers

import (
	"../factories"
	"../initializers"
	"io"
	"io/ioutil"
	"net/http"
)

func Push(w http.ResponseWriter, r *http.Request, registry *initializers.Registry) {
	rs, _ := ioutil.ReadAll(r.Body)
	s := string(rs)
	m := factories.Messages(s, r.URL.Query()["queue_id"][0])
	if m == nil {
		io.WriteString(w, "FAIL")
		return
	}

	registry.Messaging.Push(m, false)
	io.WriteString(w, "OK")
}
