package gs

import (
	"fmt"
	"testing"
)

func TestUrl2Url2SpreadsheetId(t *testing.T) {
	realId := "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms"
	url := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/edit", realId)
	if id, err := Url2SpreadsheetId(url); err != nil {
		t.Fatalf("Parse failed: %v", err)
	} else if id != realId {
		t.Fatalf("Incorrect Id %v vs. %v", id, realId)
	}
}
