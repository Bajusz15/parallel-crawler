package main

import (
	"context"
	"crypto/md5"
	"errors"
	"flag"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"parallel-crawler/crawler"
	"time"
)

var ErrUnexpectedParallelValue = fmt.Errorf("application requires parallel value greater than 0")
var addresses []string

func main() {
	logger := log.Default()
	// number of max parallel requests, by default 10
	parallel := flag.Int("parallel", 10, "Number of max parallel requests. Default is 10. (Optional)")
	flag.Parse()

	if *parallel < 1 {
		logger.Println(ErrUnexpectedParallelValue)
		os.Exit(1)
	}
	// parse addresses
	addresses = flag.Args()
	crw := crawler.NewService(logger)

	err := parallelCrawl(*parallel, logger, crw)
	if err != nil {
		logger.Fatal(err)
	}
}

func parallelCrawl(maxParallel int, logger *log.Logger, crw crawler.Service) error {
	// initialize timeout context and cancellation func
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	group, ctx := errgroup.WithContext(timeoutCtx)
	responseChan := make(chan *crawler.Response)

	// if context exceeds the deadline then stop execution
	group.Go(func() error {
		<-ctx.Done()
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return ctx.Err()
		}
		return nil
	})

	// send out all requests (fan-out)
	group.Go(func() error {
		return getResponses(ctx, cancel, crw, responseChan, logger, maxParallel)
	})

	// read responses from channel (fan-in)
	group.Go(func() error {
		for resp := range responseChan {
			logger.Printf("%s %x", resp.URL, md5.Sum(resp.Body))
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		return err
	}
	return nil
}

func getResponses(
	ctx context.Context,
	cancel context.CancelFunc,
	crw crawler.Service,
	ch chan *crawler.Response,
	logger *log.Logger,
	maxParallel int,
) error {
	defer close(ch)
	defer cancel()

	var group errgroup.Group
	group.SetLimit(maxParallel)
	for i := range addresses {
		// copy of address that's safe to use concurrently
		adr := addresses[i]
		group.Go(func() error {
			resp, err := crw.CrawlURL(ctx, adr)
			if err != nil {
				logger.Printf("error loading address %s", adr)
				return err
			}
			ch <- resp
			return nil
		})
	}

	return group.Wait()
}
