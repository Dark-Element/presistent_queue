package main

import (
	"net/http"
	"presistentQueue/handlers"
	"presistentQueue/initializers"
)

func defineRoutes(registry *initializers.Registry) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/is_alive", handlerWrapper(handlers.IsAlive, registry))
	mux.HandleFunc("/push", handlerWrapper(handlers.Push, registry))
	return mux
}

func handlerWrapper(handler func(w http.ResponseWriter, r *http.Request, reg *initializers.Registry), registry *initializers.Registry) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, registry)
	})
}
