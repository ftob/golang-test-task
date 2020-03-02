package scrapper

import (
	"context"
	gtt "github.com/ftob/golang-test-task"
	"github.com/valyala/fasthttp"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type Scrapper interface {
	Scrap(ur *url.URL) (latency int64, status int, err error)
}

type httpScrapper struct {
	client fasthttp.Client
}

func newHTTPScrapper(client fasthttp.Client) Scrapper {
	sc := httpScrapper{client: client}
	return sc
}

func (sc httpScrapper) Scrap(ur *url.URL) (latency int64, status int, err error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.SetRequestURI(ur.String())
	req.Header.SetMethod("HEAD")
	// start time
	t := time.Now().Unix()
	err = sc.client.Do(req, resp)
	t = time.Now().Unix() - t
	return t, resp.StatusCode(), err
}

type Service interface {
	StartScrapWithContext(ctx context.Context) error
}

type service struct {
	repo                  gtt.Repository
	scrappers             []Scrapper
	countScrappers        int
	currentScrapperCursor int32
}

func NewHttpService(repo gtt.Repository) Service {
	s := &service{repo: repo}

	for i := 0; i < s.repo.CountSites(); i++ {
		client := fasthttp.Client{}
		s.scrappers = append(s.scrappers, newHTTPScrapper(client))
	}
	s.currentScrapperCursor = int32(len(s.scrappers))
	return s
}

func (s *service) StartScrapWithContext(ctx context.Context) (err error) {
	var (
		ch    = make(chan string, 1)
		sites = s.repo.GetAllSites()
		wg    = sync.WaitGroup{}
	)

	defer close(ch)

	for _, site := range sites {
		ch <- site
	}
	// count wait elements
	waitElements := 0

	select {
	case site := <-ch:
		waitElements++
		go s.scrapSite(&wg, site)
	//
	case <-ctx.Done():
		wg.Add(-waitElements)
		return
	}

	wg.Wait()
	return err
}

func (s *service) scrapSite(w *sync.WaitGroup, site string) {
	w.Add(1)
	defer w.Done()
	d, err := s.repo.GetByDomain(site)
	// If domain not found must be exist a function
	if err != nil {
		return
	}
	// get first scrapper by cursor
	l, ss, err := s.nextScrapper().Scrap(d.DomainName)
	if err != nil {
		ss = 0
	}
	// store statistic information
	_ = s.repo.Store(d.DomainName.String(), l, ss)

}

func (s *service) nextScrapper() (sc Scrapper) {
	if s.currentScrapperCursor == int32(s.countScrappers-1) {
		s.currentScrapperCursor = 0
	}

	if s.currentScrapperCursor == 0 {
		sc = s.scrappers[0]
	}
	if sc == nil {
		sc = s.scrappers[s.currentScrapperCursor]
	}
	atomic.AddInt32(&s.currentScrapperCursor, 1)
	return sc
}
