package scrapper

import (
	"context"
	"github.com/ftob/golang-test-task"
	"github.com/ftob/golang-test-task/mem"
	"github.com/valyala/fasthttp"
	"net/url"
	"reflect"
	"sync"
	"testing"
)

func TestNewHttpService(t *testing.T) {
	type args struct {
		repo golang_test_task.Repository
	}
	scrappers := make([]Scrapper, 1)
	scrappers = append(scrappers, newHTTPScrapper(fasthttp.Client{}))
	tests := []struct {
		name string
		args args
		want Service
	}{
		struct {
			name string
			args args
			want Service
		}{name: "normal test", args: args{repo:mem.New([]string{})}, want: &service{
			repo:                  mem.New([]string{}),
			scrappers:             scrappers,
			countScrappers:        1,
			currentScrapperCursor: 1,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHttpService(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHttpService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpScrapper_Scrap(t *testing.T) {
	type fields struct {
		client fasthttp.Client
	}
	type args struct {
		ur *url.URL
	}
	u, _ := url.Parse("http://google.com")
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantLatency int64
		wantStatus  int
		wantErr     bool
	}{
		struct {
			name        string
			fields      fields
			args        args
			wantLatency int64
			wantStatus  int
			wantErr     bool
		}{name: "normal", fields: fields{client: fasthttp.Client{}}, args: args{ur: u}, wantLatency: 0, wantStatus: 200, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := httpScrapper{
				client: tt.fields.client,
			}
			gotLatency, gotStatus, err := sc.Scrap(tt.args.ur)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scrap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLatency != tt.wantLatency {
				t.Errorf("Scrap() gotLatency = %v, want %v", gotLatency, tt.wantLatency)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("Scrap() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}

func Test_newHTTPScrapper(t *testing.T) {
	type args struct {
		client fasthttp.Client
	}
	tests := []struct {
		name string
		args args
		want Scrapper
	}{
		struct {
			name string
			args args
			want Scrapper
		}{name: "test", args: args{client:fasthttp.Client{}}, want: &httpScrapper{client:fasthttp.Client{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newHTTPScrapper(tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newHTTPScrapper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_StartScrapWithContext(t *testing.T) {
	type fields struct {
		repo                  golang_test_task.Repository
		scrappers             []Scrapper
		countScrappers        int
		currentScrapperCursor int32
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				repo:                  tt.fields.repo,
				scrappers:             tt.fields.scrappers,
				countScrappers:        tt.fields.countScrappers,
				currentScrapperCursor: tt.fields.currentScrapperCursor,
			}
			if err := s.StartScrapWithContext(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("StartScrapWithContext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_nextScrapper(t *testing.T) {
	type fields struct {
		repo                  golang_test_task.Repository
		scrappers             []Scrapper
		countScrappers        int
		currentScrapperCursor int32
	}
	tests := []struct {
		name   string
		fields fields
		wantSc Scrapper
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				repo:                  tt.fields.repo,
				scrappers:             tt.fields.scrappers,
				countScrappers:        tt.fields.countScrappers,
				currentScrapperCursor: tt.fields.currentScrapperCursor,
			}
			if gotSc := s.nextScrapper(); !reflect.DeepEqual(gotSc, tt.wantSc) {
				t.Errorf("nextScrapper() = %v, want %v", gotSc, tt.wantSc)
			}
		})
	}
}

func Test_service_scrapSite(t *testing.T) {
	type fields struct {
		repo                  golang_test_task.Repository
		scrappers             []Scrapper
		countScrappers        int
		currentScrapperCursor int32
	}
	type args struct {
		w    *sync.WaitGroup
		site string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				repo:                  tt.fields.repo,
				scrappers:             tt.fields.scrappers,
				countScrappers:        tt.fields.countScrappers,
				currentScrapperCursor: tt.fields.currentScrapperCursor,
			}
		})
	}
}