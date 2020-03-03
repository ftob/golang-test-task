package mem

import (
	"errors"
	"fmt"
	gtt "github.com/ftob/golang-test-task"
	"net/url"
	"sync"
	"sync/atomic"
)

const UndefinedDomain = "domain name not found (please see sites.txt)"

//
type repository struct {
	sites      []string
	store      map[string]gtt.Site
	mx         sync.Mutex
	maxLatency gtt.Site
	minLatency gtt.Site
}

func New(sites []string) gtt.Repository {
	r := &repository{}

	r.sites = sites
	// Make store
	r.store = make(map[string]gtt.Site)

	for _, s := range sites {
		d, err := url.Parse(s)
		if err != nil {
			continue
		}
		r.store[s] = gtt.Site{
			DomainName: d,
			CountRequest: 0,
			Latency:      0,
			HttpStatus:   0,
		}
	}

	r.mx = sync.Mutex{}

	return r
}

//
func (r *repository) Store(domain string, latency int64, code int) (err error) {

	r.mx.Lock()
	d, err := url.Parse(domain)

	if err != nil {
		return fmt.Errorf("parse domain name (%s) error - %v", domain, err)
	}

	if s, ok := r.store[domain]; ok {
		s.Latency = latency
		s.HttpStatus = code
		s.DomainName = d
		atomic.AddUint64(&s.CountRequest, 1)

		r.store[domain] = s
		// trigger event on store
		go func() {
			_ = r.onStore(s)
		}()

	} else {
		err = errors.New(UndefinedDomain)
	}

	r.mx.Unlock()

	return err
}

//
func (r *repository) GetMaxLatency() (l gtt.Site, err error) {
	return r.maxLatency, nil
}

//
func (r *repository) GetMinLatency() (l gtt.Site, err error) {
	return r.minLatency, nil
}

//
func (r *repository) GetByDomain(domain string) (site gtt.Site, err error) {
	if s, ok := r.store[domain]; ok {
		return s, nil
	} else {
		return s, errors.New(UndefinedDomain)
	}
}

//
func (r *repository) GetAllSites() []string {
	return r.sites
}

//
func (r *repository) onStore(site gtt.Site) error {
	return r.reIndex(site)
}

// need lock
func (r *repository) reIndex(site gtt.Site) error {
	if r.maxLatency.Latency < site.Latency || r.maxLatency.Latency == 0 {
		r.maxLatency = site
	}

	if r.minLatency.Latency > site.Latency || r.maxLatency.Latency == 0 {
		r.minLatency.Latency = site.Latency
	}
	return nil
}

func (r *repository) CountSites() int {
	return len(r.sites)
}
