package main

import (
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func BOMReader(r io.Reader) io.Reader {
	return transform.NewReader(r, unicode.BOMOverride(unicode.UTF8.NewDecoder()))
}

func deferClose(err *error, f func() error) {
	newErr := f()
	if *err == nil {
		*err = newErr
	}
}

func readFrom(name string) (rc io.ReadCloser, err error) {
	if strings.HasPrefix(name, "http") {
		rsp, err := http.Get(name)
		return rsp.Body, err
	}

	f, err := os.Open(name)
	return f, err
}
