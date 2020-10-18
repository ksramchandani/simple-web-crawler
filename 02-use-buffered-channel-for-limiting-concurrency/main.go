package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	wg              = sync.WaitGroup{}
	totalGoRoutines = make(chan bool, 5)
)

type Job struct {
	url   string
	depth int
}

func crawl(c Crawler, current Job) {

	defer wg.Done()

	totalGoRoutines <- true
	defer func() {
		time.Sleep(2 * time.Second)
		<-totalGoRoutines
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
