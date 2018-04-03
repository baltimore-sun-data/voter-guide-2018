package main

import (
	"io"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func BOMReader(r io.Reader) io.Reader {
	return transform.NewReader(r, unicode.BOMOverride(unicode.UTF8.NewDecoder()))
}
