package server

import (
	"fmt"
	"github.com/fasthttp/router"
	gtt "github.com/ftob/golang-test-task"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// server and routing
type Server struct {
	router *router.Router
}

// New server and make routing map
func New(repo gtt.Repository) *Server {
	r := router.New()
	// not found page handler
	nfh := &notFoundHandler{}
	r.NotFound = nfh.handler

	// site handler
	sh := &sitesHandler{repo: repo}

	fieldKeys := []string{"method"}

	sih := sitesInstrumentingHandler{
		next: sh,
		requestCount: kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "scrapper_service",
			Name:      "request_count",
			Help:      "Number of statistics requests received.",
		}, fieldKeys),
	}
	r.GET("/sites", sih.handler)

	fhp := fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
	r.GET("/metrics", fhp)
	// index handler
	ih := &indexHandler{}
	r.GET("/", ih.handler)

	srv := &Server{
		router: r,
	}
	return srv
}

// Listen request
func (s *Server) ListenAndServe(port uint) error {
	addr := fmt.Sprintf(":%d", port)

	return fasthttp.ListenAndServe(addr, s.router.Handler)
}
