package main

import (
	"fmt"
	"time"

	"github.com/oldcookie/go-edgar"
)

func main() {
	ch := make(chan *edgar.FilingSummary)
	feed, err := edgar.NewFilingsFeed(edgar.FeedLatestFilings, ch, 10*time.Minute)
	if err != nil {
		panic(err)
	}
	feed.Monitor()

	i := 0
	ticker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case fs, ok := <-ch:
			if !ok {
				fmt.Printf("Channel closed, err %v\n", feed.Error)
				return
			}
			fmt.Printf("Received\n---------------------------------------------\n%+v\n", fs)
		case <-ticker.C:
			i++
			if i > 3 {
				feed.Close()
				return
			}
		}
	}
}
