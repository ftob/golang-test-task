package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ftob/golang-test-task/mem"
	"github.com/ftob/golang-test-task/scrapper"
	"github.com/ftob/golang-test-task/server"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	// Default enable debug information
	debug = flag.Bool("debug", true, "Enable or disable debug mod. Default --debug=true")
	// configuration HTTP server
	httpPort = flag.Uint("http.port", 8081, "setup HTTP server port. Default --http.port=8081")
)

func main() {
	var logger *zap.Logger

	// Parse parameters
	flag.Parse()

	if debug == nil {
		panic(" debug flag nil pointer exception")
	}

	if *debug {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}

	repo := mem.New([]string{"http://google.com"})
	srv := scrapper.NewHttpService(repo)
	// make context and defer cancel
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()

	go func() {
		c := time.Tick(time.Second)
		for now := range c {
			logger.Sugar().Infow("start scrap sites", "now", now)
			_ = srv.StartScrapWithContext(ctx)
			logger.Sugar().Infow("finish scrap sites", "now", time.Since(now))

		}
	}()
	serve := server.New(repo, logger)

	errs := make(chan error, 2)

	logger.Info("start HTTP server")
	errs <- serve.ListenAndServe(*httpPort)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Error("server terminated")
}
