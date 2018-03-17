package handlers

import (
	"persistentQueue/initializers"
	"github.com/valyala/fasthttp"
	"bufio"
	"io"
	"strconv"
)

func Pop(ctx *fasthttp.RequestCtx, registry *initializers.Registry) {
	topic := string(ctx.QueryArgs().Peek("topic_id"))
	targetCount, _ := strconv.ParseInt(string(ctx.QueryArgs().Peek("target_count")), 10, 64)
	targetSize, _ := strconv.ParseInt(string(ctx.QueryArgs().Peek("target_size")), 10, 64)
	b := registry.Messaging.Pop(topic, targetCount, targetSize)
	ctx.SetBodyStreamWriter(func(writer *bufio.Writer) {
		for {
			buffer := make([]byte, 1024)
			n, err := b.Read(buffer)
			if err != nil && err == io.EOF {
				writer.Flush()
				break
			}
			if n == 0 {
				break
			}
			writer.Write(buffer)
		}
	})
}
