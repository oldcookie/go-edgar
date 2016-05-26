package edgar

import (
  "time"
  "testing"
  "strings"

  "github.com/stretchr/testify/assert"
)

const teststring = "03 ENTERTAINMENT GROUP INC:0001322914:\n033 ASSET MANAGEMENT LLC /:0001114831:"

func TestGetCIKList(t *testing.T) {
  ciks := []string{"0001322914", "0001114831"}
  names := []string{"03 ENTERTAINMENT GROUP INC", "033 ASSET MANAGEMENT LLC /"}

  if testing.Short() {
    t.Skip("Skipping actual download and parsing of CIK file list in shor mode")
  } else {
    noCallHandlerFn := func (e *CIKEntry, l int) error {
      t.Errorf("This handler shouldn't be called")
      return nil
    }

    code, err := GetCIKList(noCallHandlerFn, time.Now())
    if code != 304 {
      t.Errorf("Status code should be 304 but got: %d, %s", code, err)
    }

    handlerFn := func (e *CIKEntry, l int) error {
      if e.cik != ciks[l] || e.name != names[l] {
        t.Errorf("Results didn't match - Line %d: %s, %s - Got %s, %s",
          l, ciks[l], names[l], e.cik, e.name )
      }
      return nil
    }
    parseCIKList(strings.NewReader(teststring), handlerFn)
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
