package edgar

import (
	"github.com/SlyMarbo/rss"
)

const latestFilingsURL = "https://www.sec.gov/cgi-bin/browse-edgar?action=getcurrent&type=&company=&dateb=&owner=include&start=0&count=40&output=atom"

type RSSEntryHandler func() error

type LatestFilingsFeed struct {
	feed *rss.Feed
}

func NewLatestFilingsFeed() error {
	lf := &LatestFilingsFeed{}
	feed, err := rss.Fetch(latestFilingsURL)
	if err != nil {
		return err
	}
	lf.feed = feed

	return nil
}
