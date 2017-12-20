package middlewares

import (

	"net/http"
)

func Logging(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("request started")
		//startTime := time.Now()
		next.ServeHTTP(w, r)
		//elapsedTime := time.Now().Sub(startTime)
		//fmt.Printf("time took: %s\n", elapsedTime)
		//fmt.Printf("------------------\n")
	})

}
