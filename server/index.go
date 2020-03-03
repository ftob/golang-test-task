package server

import (

	"github.com/valyala/fasthttp"
)

type indexHandler struct{}

func (ih *indexHandler) handler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("text/json")
	ctx.SetBody([]byte("{}"))
}