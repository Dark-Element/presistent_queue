package main

import (
	"net/http"
	"presistentQueue/middlewares"
	"presistentQueue/initializers"
)

func main(){
	registry := initializers.GetRegistry()
	routes := defineRoutes(registry)
	middleware := middlewares.NewMiddlares(middlewares.Logging, middlewares.Auth)
	http.ListenAndServe(":8000", middleware.Then(routes))
}