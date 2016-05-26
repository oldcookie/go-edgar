package edgar

import (
  "io"
  "time"
  "bufio"
  "strings"
  "net/http"

  "github.com/golang/glog"
)

const cikFullListURL = "http://www.sec.gov/edgar/NYU/cik.coleft.c"
const getCompanyUrl = "https://www.sec.gov/cgi-bin/browse-edgar?CIK={CIK}&owner=exclude&action=getcompany&count=10&output=atom"
const cikPH = "{CIK}"


type CIKEntry struct {
  cik, name string
}

/*
Handler for parsed CIK entry. Each time an entry is parsed,
this handler is called with an CIKEntry, and the line number from the file.
*/
type CIKEntryHandler func(*CIKEntry, int) error

/*
Retrieve the full list of CIK and Company names from EDGAR.

Takes a handler function(CIKEntryHandler), an a timestamp as arguments. Return the response
code if successful.

The timestamp, ts, is used for If-Modified-Since check so that the download would
only occur if ts < Last-Modified.
*/
func GetCIKList(fn CIKEntryHandler, ts time.Time) (int, error) {
  client := &http.Client{}
  req, err := http.NewRequest("GET", cikFullListURL, nil)
  if err != nil {
    return 0, err
  }
  if !ts.IsZero() {
    req.Header.Add("If-Modified-Since", ts.Format(time.RFC1123))
  }
  resp, err := client.Do(req)
  if err != nil {
    return 0, err
  }
  defer resp.Body.Close()

  parseCIKList(resp.Body, fn)
  return resp.StatusCode, nil
}

func parseCIKList(r io.Reader, fn CIKEntryHandler) {
  scanner := bufio.NewScanner(r)
  l := 0
  for scanner.Scan() {
    s := strings.Split(scanner.Text(), ":")
    glog.Infof("CIK list line %d: %s, %s",l, s[0], s[1])
    go fn(&CIKEntry{s[1], s[0]}, l)
    l++
  }
}

func GetCompany(cik string) (*CompanyResult, error) {
  resp, err := http.Get(strings.Replace(getCompanyUrl, cikPH, cik, -1))
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()
  c, err := parseGetCompany(resp.Body)
  glog.Info(c)
  return c, nil
}
