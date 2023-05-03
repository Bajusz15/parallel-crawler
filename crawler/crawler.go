package crawler

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Service interface {
	CrawlURL(ctx context.Context, url string) (*Response, error)
}

type Response struct {
	URL  string
	Body []byte
}

type service struct {
	logger *log.Logger
}

func NewService(l *log.Logger) Service {
	return &service{logger: l}
}

func (s *service) CrawlURL(ctx context.Context, url string) (*Response, error) {
	url = withHTTP(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		URL:  url,
		Body: body,
	}, nil
}

func withHTTP(s string) string {
	if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
		s = "https://" + s
	}
	return s
}
