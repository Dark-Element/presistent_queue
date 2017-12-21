package handlers

import (
	"persistentQueue/initializers"
	"io"
	"github.com/valyala/fasthttp"
)

func Pop(ctx *fasthttp.RequestCtx, registry *initializers.Registry) {
	b := registry.Messaging.Pop(string(ctx.QueryArgs().Peek("queue_id")), 500)
	io.WriteString(ctx, b.String())
}
