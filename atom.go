package edgar

import (
	"encoding/xml"
	"io"

	"golang.org/x/net/html/charset"
)

func parseAtom(r io.Reader) (*atomFeed, error) {
	feed := atomFeed{}
	p := xml.NewDecoder(r)
	p.CharsetReader = charset.NewReaderLabel
	err := p.Decode(&feed)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}

type atomFeed struct {
	XMLName     xml.Name        `xml:"feed"`
	Title       string          `xml:"title"`
	Description string          `xml:"subtitle"`
	Link        []atomLink      `xml:"link"`
	Items       []atomItem      `xml:"entry"`
	Updated     string          `xml:"updated"`
	CompanyInfo atomCompanyInfo `xml:"company-info"`
}

type atomItem struct {
	XMLName  xml.Name     `xml:"entry"`
	Title    string       `xml:"title"`
	Summary  string       `xml:"summary"`
	Content  string       `xml:"content"`
	Links    []atomLink   `xml:"link"`
	Date     string       `xml:"updated"`
	Category atomCategory `xml:"category"`
	ID       string       `xml:"id"`
}

type atomCategory struct {
	Label string `xml:"label,attr"`
	Term  string `xml:"term,attr"`
}

type atomLink struct {
	Href   string `xml:"href,attr"`
	Rel    string `xml:"rel,attr"`
	Type   string `xml:"type,attr"`
	Length int    `xml:"length,attr"`
}

type atomCompanyInfo struct {
	Name           string `xml:"conformed-name"`
	FiscalYearEnd  string `xml:"fiscal-year-end"`
	CIKEntry       string `xml:"cik"`
	Sic            string `xml:"assigned-sic"`
	SicDescription string `xml:"assigned-sic-desc"`
}
