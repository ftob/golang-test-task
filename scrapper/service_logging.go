package scrapper

import (
	"context"
	"go.uber.org/zap"
	"net/url"
)

type scrapperLogging struct{
	next Scrapper
	logger *zap.Logger
}

func (sl *scrapperLogging)	Scrap(ur *url.URL) (latency int64, status int, err error) {
	return sl.next.Scrap(ur)
}



type loggingService struct{
	next Service
}

func (lg *loggingService) StartScrapWithContext(ctx context.Context) error {
	return lg.next.StartScrapWithContext(ctx)
}