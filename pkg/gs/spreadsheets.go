package gs

import (
	"github.com/pkg/errors"
	"strings"
)

// Given some uri, return spreadseet id
func Url2SpreadsheetId(url string) (spreadId string, err error) {
	idLen := 44
	if len(url) < idLen {
		return "", errors.New("invalid uri")
	}
	if len(url) == idLen {
		return url, nil
	}
	errorParseFailed := errors.New("parse failed")

	checkAngGet := func(prefix string) (spreadId string, err error) {
		lenPrefix := len(prefix)
		if strings.HasPrefix(url, prefix) && len(url) >= (lenPrefix+idLen) {
			return url[lenPrefix : lenPrefix+idLen], nil
		}
		return "", errorParseFailed
	}

	pl := []string{
		"docs.google.com/spreadsheets/d/",
		"http://docs.google.com/spreadsheets/d/",
		"https://docs.google.com/spreadsheets/d/",
	}

	for _, p := range pl {
		if i, e := checkAngGet(p); e == nil {
			return i, e
		}
	}
	return "", errorParseFailed
}
