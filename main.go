package main

import (
	"persistentQueue/handlers"
	"persistentQueue/initializers"

	"github.com/valyala/fasthttp"

	"persistentQueue/middlewares"
	"os"
	"fmt"
	"syscall"
	"os/signal"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	registry := initializers.GetRegistry(sigs, done)
	routes := defineRoutes(registry)
	middleware := middlewares.NewMiddlwares(middlewares.Logging)
	go fasthttp.ListenAndServe(":8000", middleware.Then(routes))

	<-done
	fmt.Println("exiting")
}

func defineRoutes(registry *initializers.Registry) fasthttp.RequestHandler {
	m := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/push":
			handlerWrapper(handlers.Push, registry)(ctx)
		case "/pop":
			handlerWrapper(handlers.Pop, registry)(ctx)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}
	return m
}

func handlerWrapper(handler func(ctx *fasthttp.RequestCtx, registry *initializers.Registry), registry *initializers.Registry) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx){
		handler(ctx, registry)
	}
}
