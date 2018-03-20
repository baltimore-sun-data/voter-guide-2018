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
	"time"
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
	dataHeader := append([]string{}, fields...)

	for {
		fields, err = cr.Read()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return
		}

		datum := makeDatum(dataHeader, fields)
		err = saveDatum(datum)
		if err != nil {
			return err
		}
	}
}

func makeDatum(dataHeader, fields []string) map[string]interface{} {
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
	n := 1
	type question = struct {
		Question  interface{} `json:"question"`
		Answer    interface{} `json:"answer"`
		Shortname interface{} `json:"shortname"`
	}
	var questions []question
	for {
		qn := fmt.Sprintf("q%d", n)
		an := fmt.Sprintf("a%d", n)
		sn := fmt.Sprintf("sn%d", n)
		if datum[qn] == nil || datum[qn] == "" {
			break
		}
		questions = append(questions, question{
			Question:  datum[qn],
			Answer:    datum[an],
			Shortname: datum[sn],
		})
		n++
	}
	if len(questions) > 0 {
		datum["questions"] = questions
	}

	if birthday, ok := datum["birthday"].(string); ok && birthday != "" {
		birthday, err := time.Parse("1-2-2006", birthday)
		if err == nil {
			age := time.Now().Year() - birthday.Year()
			if time.Now().YearDay() < birthday.YearDay() {
				age--
			}
			datum["age"] = age
		}
	}
	return datum
}

func saveDatum(datum map[string]interface{}) (err error) {
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
	defer deferClose(&err, df.Close)

	enc := json.NewEncoder(df)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	err = enc.Encode(&datum)

	return
}
