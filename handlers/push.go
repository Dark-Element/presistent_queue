package handlers

import (
	"persistentQueue/initializers"
	"github.com/valyala/fasthttp"
	"persistentQueue/models"
)

func Push(ctx *fasthttp.RequestCtx, registry *initializers.Registry) {
	m := &models.Message{
		Data: ctx.PostBody(),
		QueueId: string(ctx.QueryArgs().Peek("queue_id")),
	}
	registry.Messaging.Push(m, false)
	ctx.SetBody([]byte(""))
}
