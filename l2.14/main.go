package main

import (
	"fmt"
	"time"
)

func or(channels ...<-chan interface{}) <-chan interface{} {
	ch := make(chan interface{})
	go func(ch chan<- interface{}, channels ...<-chan interface{}) {
		defer close(ch)
		if len(channels) == 1 {
			<-channels[0]
		} else if len(channels) > 1 {
			select {
			case <-channels[len(channels)-1]:
				break
			case <-or(channels[:len(channels)-1]...):
				break
			}
		}
	}(ch, channels...)
	return ch
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(4*time.Second),
		sig(2*time.Second),
		sig(5*time.Second),
	)
	fmt.Printf("done after %v\n", time.Since(start))
}
