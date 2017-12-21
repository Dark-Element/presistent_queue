package main

import (
	"persistentQueue/handlers"
	"persistentQueue/initializers"

	"github.com/valyala/fasthttp"

	"persistentQueue/middlewares"
)

func main() {
	registry := initializers.GetRegistry()
	routes := defineRoutes(registry)
	middleware := middlewares.NewMiddlwares(middlewares.Logging)
	fasthttp.ListenAndServe(":8000", middleware.Then(routes))
	//http.ListenAndServe(":8000", )
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
