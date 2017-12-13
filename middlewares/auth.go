package middlewares

import (
	"net/http"
	"fmt"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Authorized")
		next.ServeHTTP(w, r)
	})
}

func autorizeRequest(r *http.Request) {

}
