package server

import (
	"github.com/valyala/fasthttp"
)

//
type notFoundHandler struct{}

func (nfh *notFoundHandler) handler(ctx *fasthttp.RequestCtx)  {
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.SetContentType("text/json")
	ctx.SetBody([]byte(`{"message": "Page not found"}`))
}
