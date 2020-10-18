package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	wg              = sync.WaitGroup{}
	totalGoRoutines = 5
	buffered        = make(chan bool, totalGoRoutines)
)

type Job struct {
	url   string
	depth int
}

func crawl(c Crawler, current Job) {
	defer wg.Done()

	// This is the only change from the previous method
	// We limit the number of go routines by adding a buffered channel
	// that will only allow sending upto totalGoRoutines messages on the channel
	// The reason to do the wg.Add(1) by the url is because it is almost impossible
	// to use a worker pool with channels for this scenario
	// If we used a worker pool, a worker will receive a url from a channel
	// Now, after crawling, more urls needs to be crawled
	// Sending on the same channel will deadlock it even with buffering
	// Sending it on a different channel causes additional complexity of now receiving
	// it from the other channel and co-ordinating when the code terminates

	buffered <- true
	defer func() {
		time.Sleep(2 * time.Second)
		<-buffered
	}()

	if current.depth <= 0 {
		return
	}

	newURLS := c.Fetch(current.url)

	fmt.Printf("current = %v. newurls = %v\n", current.url, newURLS)

	for _, u := range newURLS {
		wg.Add(1)
		j := Job{
			url:   u,
			depth: current.depth - 1,
		}
		go crawl(c, j)
	}
}

func main() {
	// Create start url and depth
	start := Job{
		url:   "www.google.com",
		depth: 3,
	}

	// mockCrawler implements the Crawler interface and helps us fetching mock urls
	m := mockCrawler{}

	wg.Add(1)
	go crawl(m, start)

	wg.Wait()
}
