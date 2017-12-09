package main

import (
	"net/http"
	"presistentQueue/middlewares"
)

func main(){
	routes := defineRoutes()
	middleware := middlewares.NewMiddlares(middlewares.Logging, middlewares.Auth)
	http.ListenAndServe(":8000", middleware.Then(routes))
}


