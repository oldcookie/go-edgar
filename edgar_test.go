package edgar

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const teststring = "03 ENTERTAINMENT GROUP INC:0001322914:\n033 ASSET MANAGEMENT LLC /:0001114831:"

type Handler struct {
	t *testing.T
}

var ciks = []string{"0001322914", "0001114831"}
var names = []string{"03 ENTERTAINMENT GROUP INC", "033 ASSET MANAGEMENT LLC /"}

func (h *Handler) ProcessCIKEntry(cik, name string, line int) {
	if cik != ciks[line] || name != names[line] {
		h.t.Errorf("Results didn't match - Line %d: %s, %s - Got %s, %s", line,
			ciks[line], names[line], cik, name)
	}
}

type NoCallHandler struct {
	t *testing.T
}

func (h *NoCallHandler) ProcessCIKEntry(cik, name string, line int) {
	h.t.Errorf("This handler shouldn't be called, %s, %s, %d", cik, name, line)
}

func TestGetCIKList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping actual download and parsing of CIK file list in short mode")
	} else {
		code, err := GetCIKList(&NoCallHandler{t}, time.Now())
		if code != 304 {
			t.Errorf("Status code should be 304 but got: %d, %s", code, err)
		}
		parseCIKList(strings.NewReader(teststring), &Handler{t})
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
