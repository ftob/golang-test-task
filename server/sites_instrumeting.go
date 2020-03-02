package server

import (
	"fmt"
	"github.com/go-kit/kit/metrics"
	"github.com/valyala/fasthttp"
)

const (
	MethodMaxLatency = "max"
	MethodMinLatency = "min"
	MethodDomain     = "domain"
)

//
type sitesInstrumentingHandler struct {
	next         *sitesHandler
	requestCount metrics.Counter
}

// Prometheus
func (sih *sitesInstrumentingHandler) handler(ctx *fasthttp.RequestCtx) {
	defer func() {
		var method string

		method = "non"

		if p := ctx.QueryArgs().Peek("max"); len(p) > 0 {
			method = MethodMaxLatency
		}

		if p := ctx.QueryArgs().Peek("max"); len(p) > 0{
			method = MethodMinLatency
		}

		if p := ctx.QueryArgs().Peek("max"); len(p) > 0{
			method = MethodDomain
			sih.requestCount.With("domain", fmt.Sprintf("%s",p))
		}
		sih.requestCount.With("method", method).Add(1)
	}()

	sih.next.handler(ctx)
}
