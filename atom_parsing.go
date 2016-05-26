package edgar

import (
  "io"
  "encoding/xml"
  "time"
  "golang.org/x/net/html/charset"
)

type CompanyResult struct {
  XMLName xml.Name      `xml:"feed"`
  Info    CompanyInfo   `xml:"company-info"`
  Title   string        `xml:"title"`
  Updated time.Time     `xml:"updated"`
}

type CompanyInfo struct {
  XMLName         xml.Name          `xml:"company-info"`
  Cik             string            `xml:"cik"`
  Name            string            `xml:"conformed-name"`
  FiscalYearEnd   string            `xml:"fiscal-year-end"`
  Sic             string            `xml:"assigned-sic"`
  SicDescription  string            `xml:"assigned-sic-desc"`
  StateLocation   string            `xml:"state-location"`
  StateOfIncorp   string            `xml:"state-of-incorporation"`
  Addresses       []StreetAddress  `xml:"addresses>address"`
}

type StreetAddress struct {
  Type            string      `xml:"type,attr"`
  City            string      `xml:"city"`
  State           string      `xml:"state"`
  Street1         string      `xml:"street1"`
  Street2         string      `xml:"street2"`
  Street3         string      `xml:"street3"`
  Zip             string      `xml:"zip"`
  Phone           string      `xml:"phone"`
}

func parseGetCompany(r io.Reader) (*CompanyResult, error) {
  d := xml.NewDecoder(r)
  d.CharsetReader = charset.NewReaderLabel
  c := CompanyResult{}
  err := d.Decode(&c)
  if err != nil {
    return nil, err
  }
  return &c, err
}
