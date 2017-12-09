package middlewares

import (
	"net/http"
)

type Middlewares struct {
	myMiddleWare []func(next http.Handler) http.Handler
}

func NewMiddlares(middlewares ...func(next http.Handler) http.Handler) *Middlewares{
	chain := Middlewares{myMiddleWare: middlewares}
	return &chain
}

func (md *Middlewares) Then(handler http.Handler) http.Handler{
	maxIdx := len(md.myMiddleWare) - 1
	builtChain := handler
	for idx, _  := range md.myMiddleWare{
		builtChain = md.myMiddleWare[maxIdx - idx](builtChain)
	}
	return builtChain
}

