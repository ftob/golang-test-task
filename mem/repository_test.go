package mem

import (
	gtt "github.com/ftob/golang-test-task"
	"reflect"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		sites []string
	}
	tests := []struct {
		name string
		args args
		want gtt.Repository
	}{
		struct {
			name string
			args args
			want gtt.Repository
		}{name: "normal new", args: args{sites: []string{"google.com", "mail.ru"}}, want: &repository{
			sites:      []string{"google.com", "mail.ru"},
			store:      make(map[string]gtt.Site),
			mx:         sync.Mutex{},
			doOnce:     sync.Once{},
			maxLatency: gtt.Site{},
			minLatency: gtt.Site{},
		}},
		struct {
			name string
			args args
			want gtt.Repository
		}{name: "empty site list", args: args{sites: []string{}}, want: New([]string{})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.sites); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}

}

func Test_repository_Store(t *testing.T) {
	type fields struct {
		sites      []string
		store      map[string]gtt.Site
		mx         sync.Mutex
		doOnce     sync.Once
		maxLatency gtt.Site
		minLatency gtt.Site
	}
	type args struct {
		domain  string
		latency int64
		code    int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		struct {
			name    string
			fields  fields
			args    args
			wantErr bool
		}{name: "normal behavior", fields: struct {
			sites      []string
			store      map[string]gtt.Site
			mx         sync.Mutex
			doOnce     sync.Once
			maxLatency gtt.Site
			minLatency gtt.Site
		}{sites: []string{"test.com"}, store: make(map[string]gtt.Site), mx: sync.Mutex{}, doOnce: sync.Once{}, maxLatency: gtt.Site{
			DomainName:   nil,
			CountRequest: 1,
			Latency:      2,
			HttpStatus:   3,
		}, minLatency: gtt.Site{
			DomainName:   nil,
			CountRequest: 0,
			Latency:      0,
			HttpStatus:   0,
		}}, args: args{
			domain:  "test.com",
			latency: 1,
			code:    200,
		}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{
				sites:      tt.fields.sites,
				store:      tt.fields.store,
				mx:         tt.fields.mx,
				doOnce:     tt.fields.doOnce,
				maxLatency: tt.fields.maxLatency,
				minLatency: tt.fields.minLatency,
			}
			if err := r.Store(tt.args.domain, tt.args.latency, tt.args.code); (err != nil) != tt.wantErr {
				t.Errorf("Store() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_repository_reIndex(t *testing.T) {
	type fields struct {
		sites      []string
		store      map[string]gtt.Site
		mx         sync.Mutex
		doOnce     sync.Once
		maxLatency gtt.Site
		minLatency gtt.Site
	}
	type args struct {
		site gtt.Site
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		struct {
			name    string
			fields  fields
			args    args
			wantErr bool
		}{name: "start reindex", fields: struct {
			sites      []string
			store      map[string]gtt.Site
			mx         sync.Mutex
			doOnce     sync.Once
			maxLatency gtt.Site
			minLatency gtt.Site
		}{sites: []string{}, store: map[string]gtt.Site{}, mx: sync.Mutex{}, doOnce: sync.Once{}, maxLatency: gtt.Site{}, minLatency: gtt.Site{}}, args: args{site: gtt.Site{
			CountRequest: 0,
			Latency:      1,
			HttpStatus:   200,
		}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{
				sites:      tt.fields.sites,
				store:      tt.fields.store,
				mx:         tt.fields.mx,
				doOnce:     tt.fields.doOnce,
				maxLatency: tt.fields.maxLatency,
				minLatency: tt.fields.minLatency,
			}
			if err := r.reIndex(tt.args.site); (err != nil) != tt.wantErr {
				t.Errorf("reIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
