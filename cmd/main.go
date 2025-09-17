package main

import (
	"fmt"
	"or"
	"time"
)

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() { defer close(c); time.Sleep(after) }()
	return c
}

func main() {
	start := time.Now()
	<-or.Or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v\n", time.Since(start))
}
