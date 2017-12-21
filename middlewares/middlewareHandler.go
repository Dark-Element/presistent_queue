package middlewares

import (
	"github.com/valyala/fasthttp"
)

type Middlewares struct {
	myMiddleWare []func(next func(next *fasthttp.RequestCtx)) func(next *fasthttp.RequestCtx)
}

func NewMiddlwares(middlewares ...func(next func(next *fasthttp.RequestCtx))  func(next *fasthttp.RequestCtx)) *Middlewares {
	chain := Middlewares{myMiddleWare: middlewares}
	return &chain
}

func (md *Middlewares) Then(handler func(next *fasthttp.RequestCtx)) func(ctx *fasthttp.RequestCtx) {
	maxIdx := len(md.myMiddleWare) - 1
	builtChain := handler
	for idx, _ := range md.myMiddleWare {
		builtChain = md.myMiddleWare[maxIdx-idx](builtChain)
	}
	return builtChain
}
