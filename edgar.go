package edgar

import (
	"bufio"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang/glog"
)

const (
	cikFullListURL = "https://www.sec.gov/edgar/NYU/cik.coleft.c"
	getCompanyURL  = "https://www.sec.gov/cgi-bin/browse-edgar?CIK={CIK}&owner=exclude&action=getcompany&count=10&output=atom"
	cikPH          = "{CIK}"
)

// CIKEntryHandler - Handler for parsed CIK entry. Each time an entry is parsed,
// this handler is called with an CIKEntry, and the line number from the file.
type CIKEntryHandler interface {
	ProcessCIKEntry(cik, name string, line int)
}

// GetCIKList - Retrieve the full list of CIK and Company names from EDGAR.
// Takes a handler function(CIKEntryHandler), an a timestamp as arguments. Return the response
// code if successful.
//
// The timestamp, ts, is used for If-Modified-Since check so that the download would
// only occur if ts < Last-Modified.
func GetCIKList(h CIKEntryHandler, ts time.Time) (int, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", cikFullListURL, nil)
	if err != nil {
		return 0, err
	}
	if !ts.IsZero() {
		req.Header.Add("If-Modified-Since", ts.UTC().Format(time.RFC1123))
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	parseCIKList(resp.Body, h)
	return resp.StatusCode, nil
}

func parseCIKList(r io.Reader, h CIKEntryHandler) {
	scanner := bufio.NewScanner(r)

	for l := 0; scanner.Scan(); l++ {
		s := strings.Split(scanner.Text(), ":")
		glog.Infof("CIK list line %d: %s, %s", l, s[0], s[1])
		h.ProcessCIKEntry(s[1], s[0], l)
	}
}

// GetCompany - Get the company summary from EDGAR.
//
// Company info can be retrieved by either CIK, or the Ticker.  Either can
// be used for the cik argument.  Returns CompanryResult, nil if successful.
func GetCompany(cik string) (*CompanyResult, error) {
	resp, err := http.Get(strings.Replace(getCompanyURL, cikPH, cik, -1))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	c, _ := parseGetCompany(resp.Body)
	glog.Info(c)
	return c, nil
}
