package edgar

import (
	"fmt"
	"net/http"
	"testing"
)

func TestParseAtom(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", string(FeedLatestFilings), nil)
	if err != nil {
		t.Error(err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	defer resp.Body.Close()

	f, err := parseAtom(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(f)
}
