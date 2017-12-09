package main

import (
	"net/http"
	"presistentQueue/controllers"
	"presistentQueue/initializers"
)




func defineRoutes() *http.ServeMux{
	registry := initializers.GetRegistry() //get registry can be a overwritten per every controller
	mux := http.NewServeMux()
	mux.HandleFunc("/is_alive", controllerWrapper(controllers.IsAlive, registry) )
	return mux
}


func controllerWrapper(controller func(w http.ResponseWriter, r *http.Request, reg *initializers.Registry),
	registry *initializers.Registry) http.HandlerFunc{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		controller(w, r, registry)
	})
}