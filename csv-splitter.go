package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	files := FileList(os.Args[1:])
	if err := files.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func deferClose(err *error, f func() error) {
	newErr := f()
	if *err == nil {
		*err = newErr
	}
}

type FileList []string

var errMissingInfo = fmt.Errorf("missing filename/directory")

func (fl FileList) Execute() error {
	for _, fn := range fl {
		log.Printf("open: %s", fn)
		f, err := os.Open(fn)
		if err != nil {
			return err
		}
		err = parseAndSplit(f)
		if err != nil {
			if err == errMissingInfo {
				return fmt.Errorf("read %s: %v", fn, err)
			}
			return err
		}
	}
	return nil
}

func parseAndSplit(src io.ReadCloser) (err error) {
	defer deferClose(&err, src.Close)

	cr := csv.NewReader(src)
	cr.Comment = '#'
	cr.FieldsPerRecord = -1
	cr.ReuseRecord = true

	fields, err := cr.Read()

	// Save headers for each row of dict
	dataHeader := make(map[int]string, len(fields))
	for i, field := range fields {
		dataHeader[i] = field
	}

	for {
		fields, err = cr.Read()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return
		}

		datum := make(map[string]interface{}, len(fields))
		for i, val := range fields {
			if val == "" {
				continue
			}
			if n, err := strconv.ParseFloat(val, 64); err == nil && (val != "directory" || val != "filename") {
				datum[dataHeader[i]] = n
				continue
			}
			datum[dataHeader[i]] = val
		}

		dir, _ := datum["directory"].(string)
		fn, _ := datum["filename"].(string)

		if dir == "" || fn == "" {
			return errMissingInfo
		}

		log.Printf("creating %s/%s", dir, fn)

		if err = os.MkdirAll(dir, 0777); err != nil {
			return err
		}

		dfn := filepath.Join(dir, fn)
		df, err := os.Create(dfn)
		if err != nil {
			return err
		}

		// Can't defer close in a loop

		enc := json.NewEncoder(df)
		enc.SetIndent("", "  ")
		enc.SetEscapeHTML(false)

		if err = enc.Encode(&datum); err != nil {
			df.Close()
			return err
		}

		if err = df.Close(); err != nil {
			return err
		}
	}
}
