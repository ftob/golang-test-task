package server

import (
	"fmt"
	gtt "github.com/ftob/golang-test-task"
	"github.com/valyala/fasthttp"
)

//
type sitesHandler struct {
	repo gtt.Repository
}

//
func (sh *sitesHandler) handler(ctx *fasthttp.RequestCtx) {
	var (
		r   []byte
		err error
		d gtt.Site
	)
	ctx.SetContentType("text/json")

	// check arguments max
	if p := ctx.QueryArgs().Peek("max"); len(p) > 0 {
		d, err = sh.repo.GetMaxLatency()
		r = []byte(fmt.Sprintf(`{"result": %d, "err": "%v"}`, d.Latency, err))
	}

	if p := ctx.QueryArgs().Peek("min"); len(p) > 0 {
		d, err = sh.repo.GetMinLatency()
		r = []byte(fmt.Sprintf(`{"result": %d, "err": "%v"}`, d.Latency, err))
	}

	if p := ctx.QueryArgs().Peek("domain"); len(p) > 0 {
		d, err := sh.repo.GetByDomain(string(p))
		if err != nil {
			r = []byte(
				fmt.Sprintf(
					`{"domain": "%s", "latency": %d, "http": %d, "countRequest": %d, err: "%v"}`,
					d.DomainName.String(), d.Latency, d.HttpStatus, d.CountRequest, err))
		} else {
			r = []byte(fmt.Sprintf(`{"err": %v}`, err))
		}
	}

	if r != nil {
		ctx.SetBody(r)
		ctx.SetStatusCode(fasthttp.StatusOK)
	} else {
		ctx.SetBody([]byte(`{"status": "Bad request"}`))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)

	}

}