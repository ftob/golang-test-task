package golang_test_task

import (
	"net/url"
)

type Site struct {
	DomainName   *url.URL
	CountRequest uint64
	Latency      int64
	HttpStatus   int
}

type Repository interface {
	Store(domain string, latency int64, code int) error
	GetMaxLatency() (l Site, err error)
	GetMinLatency() (l Site, err error)
	GetByDomain(domain string) (site Site, err error)
	CountSites() int
	GetAllSites() []string
}
