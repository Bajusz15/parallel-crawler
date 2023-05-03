package crawler

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"
)

func Test_service_CrawlURL(t *testing.T) {
	s := &service{
		logger: log.Default(),
	}

	badURL := "dasdasdas"
	httpMissing := "twitter.com"
	goodURL := "https://twitter.com"
	want := "https://" + httpMissing

	_, err := s.CrawlURL(context.Background(), badURL)
	if err == nil {
		t.Errorf("expected an error but got nil")
	}

	_, err = s.CrawlURL(context.Background(), goodURL)
	if err != nil {
		t.Errorf("expected error to be nil but got %s", err)
	}
	// case when http(s) protocol is missing
	got, err := s.CrawlURL(context.Background(), httpMissing)
	if err != nil {
		t.Errorf("expected error to be nil but got %s", err)
	}
	if got.URL != want {
		t.Errorf("unexpected response when protocol is missing, got: %s, wanted: %s", got.URL, want)
	}

	// case when context deadline is exceeded
	exceeded, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	_, err = s.CrawlURL(exceeded, httpMissing)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("unexpected error, got: %s but wanted: %s", err, context.DeadlineExceeded)
	}
}
