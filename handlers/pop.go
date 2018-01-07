package handlers

import (
	"persistentQueue/initializers"
	"github.com/valyala/fasthttp"
)

func Pop(ctx *fasthttp.RequestCtx, registry *initializers.Registry) {
	b := registry.Messaging.Pop(string(ctx.QueryArgs().Peek("queue_id")), 500)
	ctx.SetBodyStream(b, -1)
}
