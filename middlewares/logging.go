package middlewares

import (
	"github.com/valyala/fasthttp"
	"fmt"
	"time"
)

func Logging(next func(ctx *fasthttp.RequestCtx)) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx){
		st := time.Now()
		next(ctx)
		fmt.Printf("Request took: %s\n",time.Now().Sub(st))
	}

}
