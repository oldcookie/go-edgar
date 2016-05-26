package edgar

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const teststring = "03 ENTERTAINMENT GROUP INC:0001322914:\n033 ASSET MANAGEMENT LLC /:0001114831:"

func TestGetCIKList(t *testing.T) {
	var ciks = []string{"0001322914", "0001114831"}
	var names = []string{"03 ENTERTAINMENT GROUP INC", "033 ASSET MANAGEMENT LLC /"}

	if testing.Short() {
		t.Skip("Skipping actual download and parsing of CIK file list in shor mode")
	} else {
		entries, err := GetCIKList(time.Now(), 100)
		if err != nil {
			t.Error(err)
		}
		for entry := range entries {
			t.Error("This should never be reached", entry)
		}

		ch := make(chan *CIKEntry, 10)
		parseCIKList(strings.NewReader(teststring), ch)
		close(ch)
		for e := range ch {
			if e.cik != ciks[e.lineNo] || e.name != names[e.lineNo] {
				t.Errorf("Results didn't match - Line %d: %s, %s - Got %s, %s",
					e.lineNo, ciks[e.lineNo], names[e.lineNo], e.cik, e.name)
			}
		}
	}
}

func TestGetCompany(t *testing.T) {
	c, err := GetCompany("AAPL")
	if err != nil {
		t.Error(err)
	}
	assert := assert.New(t)

	t.Log(c)
	assert.Equal(c.Info.Cik, "0000320193", "CIK should be the same")
	assert.Equal(c.Info.Name, "APPLE INC", "Name should be the APPLE INC")
}
