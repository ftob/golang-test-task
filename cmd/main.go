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

	repo := mem.New([]string{})
	srv := scrapper.NewHttpService(repo)
	cron := time.NewTimer(time.Minute)
	// make context and defer cancel
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()

	go func() {
		<-cron.C
		_ = srv.StartScrapWithContext(ctx)
	}()
	serve := server.New(repo)

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
