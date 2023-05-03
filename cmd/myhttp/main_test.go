package main

import (
	"context"
	"errors"
	"log"
	"parallel-crawler/crawler"
	"testing"
)

var mockError = errors.New("mock error")

type mockCrawler struct{}

func (m *mockCrawler) CrawlURL(ctx context.Context, url string) (*crawler.Response, error) {
	switch url {
	case "mock":
		return nil, mockError
	default:
		return &crawler.Response{
			URL:  url,
			Body: nil,
		}, nil
	}
}

func Test_parallelCrawl(t *testing.T) {
	logger := log.Default()
	mc := new(mockCrawler)

	tests := []struct {
		name    string
		crw     crawler.Service
		addrs   []string
		wantErr bool
	}{
		{
			name:    "no error",
			crw:     mc,
			addrs:   []string{"twitter.com"},
			wantErr: false,
		},
		{
			name:    "crawler returns error",
			crw:     mc,
			addrs:   []string{"twitter.com", "mock"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		addresses = tt.addrs
		t.Run(tt.name, func(t *testing.T) {
			err := parallelCrawl(5, logger, tt.crw)
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %s", err)
			}
		})
	}
}
