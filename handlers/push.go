package handlers

import (
	"net/http"
	"presistentQueue/initializers"
	"io"
	"presistentQueue/factories"
	"io/ioutil"
)

func Push(w http.ResponseWriter, r *http.Request, registry *initializers.Registry){
	rs,_ := ioutil.ReadAll(r.Body)
	s := string(rs)
	m := factories.Messages(s)
	if m == nil{
		io.WriteString(w, "FAIL")
		return
	}
	registry.Messaging.Push(m)
	io.WriteString(w, "OK")
}
