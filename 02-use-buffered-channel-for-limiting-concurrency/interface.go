package main

type Crawler interface {
	Fetch(url string) (urls []string)
}
