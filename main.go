package main

import (
	"./handlers"
	"./initializers"
	"./middlewares"
	"net/http"
)

func main() {
	registry := initializers.GetRegistry()
	routes := defineRoutes(registry)
	middleware := middlewares.NewMiddlares(middlewares.Logging, middlewares.Auth)
	http.ListenAndServe(":8000", middleware.Then(routes))
}

func defineRoutes(registry *initializers.Registry) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/is_alive", handlerWrapper(handlers.IsAlive, registry))
	mux.HandleFunc("/push", handlerWrapper(handlers.Push, registry))
	mux.HandleFunc("/pop", handlerWrapper(handlers.Pop, registry))
	return mux
}

func handlerWrapper(handler func(w http.ResponseWriter, r *http.Request, reg *initializers.Registry), registry *initializers.Registry) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, registry)
	})
}
