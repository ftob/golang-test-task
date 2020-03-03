package server

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type sitesLoggingHandler struct {
	next *sitesInstrumentingHandler
	logger *zap.Logger
}

// Logging wrapper
func (h *sitesLoggingHandler) handler(ctx *fasthttp.RequestCtx) {
	defer func() {
		h.logger.Sugar().Infow("params",
			"max", string(ctx.QueryArgs().Peek("max")),
			"min", string(ctx.QueryArgs().Peek("min")),
			"domain", string(ctx.QueryArgs().Peek("domain")))
	}()
	h.next.handler(ctx)
}
