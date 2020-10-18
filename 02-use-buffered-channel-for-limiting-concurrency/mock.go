package main

import "fmt"

type mockCrawler struct{}

func (c mockCrawler) Fetch(url string) []string {

	urls := []string{}

	for i := 0; i <= 10; i++ {
		urls = append(urls, fmt.Sprintf("%s/%d", url, i))
	}
	return urls

}
