package edgar

import (
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/hashicorp/golang-lru"
)

// FilingSummary - Summary entry for a filing
type FilingSummary struct {
	ID          string
	Title       string
	Summary     string
	Links       []string
	FormType    string
	Date        time.Time
	CIK         string
	CIKType     string
	AccessionNo string
}

const (
	// FeedLatestFilings - URL for getting the latest filings
	FeedLatestFilings = "https://www.sec.gov/cgi-bin/browse-edgar?action=getcurrent&type=&company=&dateb=&owner=include&start=0&count=10&output=atom"
)

var titleRE = regexp.MustCompile(`(.+)\s+\((\w+)\)\s+\((\w+)\)`)
var accessionRE = regexp.MustCompile(`.+accession-number=([\w+|-]+)`)

// FilingsFeed - feed for latest filings
type FilingsFeed struct {
	client       *http.Client
	URL          string
	UpdatePeriod time.Duration
	done         chan bool
	ch           chan *FilingSummary
	recent       *lru.Cache
	Error        error
}

// NewFilingsFeed - Create a new feed instance to monitor a remote feed
func NewFilingsFeed(URL string, ch chan *FilingSummary, updatePeriod time.Duration) (*FilingsFeed, error) {
	client := &http.Client{}
	done := make(chan bool)

	recent, err := lru.New(1000) // keep track of 1000 entries
	if err != nil {
		return nil, err
	}
	f := &FilingsFeed{client, URL, updatePeriod, done, ch, recent, nil}
	return f, nil
}

// Monitor - updates for remote feed
func (f *FilingsFeed) Monitor() {
	go func() {
		ticker := time.NewTicker(f.UpdatePeriod)
		f.update(false)
		for {
			select {
			case <-ticker.C:
				f.update(true)
			case <-f.done:
				ticker.Stop()
				return
			}
		}
	}()
}

// Close the feed
func (f *FilingsFeed) Close() {
	close(f.ch)
	close(f.done)
}

func (f *FilingsFeed) fetch(start int) (*atomFeed, error) {
	req, err := http.NewRequest("GET", f.URL, nil)
	if err != nil {
		return nil, err
	}

	if start != 0 {
		v := req.URL.Query()
		v.Set("start", strconv.Itoa(start))
		req.URL.RawQuery = v.Encode()
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	feed, err := parseAtom(resp.Body)
	if err != nil {
		return nil, err
	}
	return feed, nil
}

func (f *FilingsFeed) error(err error) {
	glog.Errorf("Error encountered %v\n", err)
	f.Error = err
	f.Close()
}

func (f *FilingsFeed) update(getAllUnseen bool) {
	for start := 0; ; {
		feed, err := f.fetch(start)
		if err != nil {
			f.error(err)
			return
		}

		for _, item := range feed.Items {
			if seen, _ := f.recent.ContainsOrAdd(item.Title+item.ID, true); !seen {
				date, err := time.Parse("2006-01-02T15:04:05-07:00", item.Date)
				if err != nil {
					f.error(err)
					return
				}

				var title, CIK, CIKType, accessionNo string
				if matches := titleRE.FindStringSubmatch(item.Title); len(matches) > 0 {
					title, CIK, CIKType = matches[1], matches[2], matches[3]
				} else {
					title = item.Title
				}
				if matches := accessionRE.FindStringSubmatch(item.ID); len(matches) > 0 {
					accessionNo = matches[1]
				}

				fs := FilingSummary{
					ID:          item.ID,
					Title:       title,
					Summary:     item.Summary,
					FormType:    item.Category.Term,
					Links:       make([]string, len(item.Links)),
					Date:        date,
					CIK:         CIK,
					CIKType:     CIKType,
					AccessionNo: accessionNo,
				}

				for i, l := range item.Links {
					fs.Links[i] = l.Href
				}
				f.ch <- &fs
			} else {
				return
			}
		}
		if !getAllUnseen {
			return
		}
		start += len(feed.Items)
	}
}
