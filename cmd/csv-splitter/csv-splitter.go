package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
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

func (fl FileList) Execute() error {
	for _, fn := range fl {
		log.Printf("open: %s", fn)
		f, err := os.Open(fn)
		if err != nil {
			return err
		}
		err = parseAndSplit(f)
		if err != nil {
			return fmt.Errorf("processing %s: %v", fn, err)
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

		var datum map[string]interface{}
		datum, err = makeDatum(dataHeader, fields)
		if err != nil {
			return err
		}
		err = saveDatum(datum)
		if err != nil {
			return err
		}
	}
}

// Checks that text is ASCII or a whitelisted Unicode character
var unicodeChecker = regexp.MustCompile(`[^\x00-\x7F‘’–—•”“ñ§½…¢óï]`)

func makeDatum(dataHeader, fields []string) (map[string]interface{}, error) {
	datum := make(map[string]interface{}, len(fields))
	for i, val := range fields {
		if val == "" {
			continue
		}
		if unicodeChecker.MatchString(val) {
			return nil, fmt.Errorf("file has Unicode errors: %q", val)
		}
		datum[dataHeader[i]] = normalize(val)
	}
	n := 1
	type question = struct {
		Question  string `json:"question"`
		Answer    string `json:"answer"`
		Shortname string `json:"shortname"`
	}
	var questions []question
	if get(datum, "a1") != "" {
		for {
			qn := fmt.Sprintf("q%d", n)
			an := fmt.Sprintf("a%d", n)
			sn := fmt.Sprintf("sn%d", n)
			if get(datum, qn) == "" {
				break
			}
			questions = append(questions, question{
				Question:  get(datum, qn),
				Answer:    get(datum, an),
				Shortname: get(datum, sn),
			})
			// Clean up JSON output
			delete(datum, qn)
			delete(datum, an)
			delete(datum, sn)
			n++
		}
	}
	if len(questions) > 0 {
		datum["questions"] = questions
	}

	if dob := get(datum, "dob"); dob != "" && get(datum, "age") == "" {
		electionDay, _ := time.Parse("1/2/2006", "11/6/2018")
		birthdate, err := time.Parse("1/2/2006", dob)
		if err == nil {
			age := electionDay.Year() - birthdate.Year()
			if electionDay.YearDay() < birthdate.YearDay() {
				age--
			}
			if age < 10 || age > 100 {
				return nil, fmt.Errorf("bad birthday: %s", dob)
			}
			datum["age"] = age
		} else {
			return nil, fmt.Errorf("could not parse dob: %v", err)
		}
	}

	dir := get(datum, "directory")
	fn := get(datum, "filename")
	if dir == "" && fn == "" {
		if race := get(datum, "race"); race != "" {
			if district := get(datum, "district"); district != "" {
				datum["directory"] = fmt.Sprintf("content/%s/district-%s",
					race, district)
			} else {
				datum["directory"] = fmt.Sprintf("content/%s", race)
			}
		} else if mun := get(datum, "race-municipality"); mun != "" {
			mun = slugify(mun)
			name := slugify(get(datum, "race-name"))
			district := slugify(get(datum, "race-district"))
			datum["directory"] = fmt.Sprintf("content/%s-county/%s/district-%v",
				mun, name, district)
		}

		fullname := get(datum, "full-name")
		datum["filename"] = fmt.Sprintf("%s.md", slugify(fullname))
	}

	if title := get(datum, "title"); title == "" {
		fullname := get(datum, "full-name")
		datum["title"] = fullname
	}

	return datum, nil
}

var strictNormalizedForms = map[string]interface{}{
	"NO":  false,
	"YES": true,
}

var normalizedForms = map[string]interface{}{
	"democrat":    "Democrat",
	"democratic":  "Democrat",
	"independent": "Independent",
	"republican":  "Republican",
}

func normalize(s string) interface{} {
	s = strings.TrimSpace(s)

	// Social media normalization
	if u, err := url.Parse(s); err == nil &&
		(strings.HasSuffix(u.Hostname(), "facebook.com") ||
			strings.HasSuffix(u.Hostname(), "twitter.com") ||
			strings.HasSuffix(u.Hostname(), "instagram.com")) {
		return strings.Trim(u.Path, "/")
	}
	if n, err := strconv.ParseFloat(s, 64); err == nil {
		return n
	}
	if n, ok := strictNormalizedForms[s]; ok {
		return n
	}
	if n, ok := normalizedForms[strings.ToLower(s)]; ok {
		return n
	}

	return s
}

var errMissingInfo = fmt.Errorf("missing filename/directory")

func saveDatum(datum map[string]interface{}) (err error) {
	dir := get(datum, "directory")
	fn := get(datum, "filename")

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

func get(m map[string]interface{}, key string) string {
	s, _ := m[key].(string)
	if f, ok := m[key].(float64); ok {
		s = fmt.Sprintf("%v", f)
	}
	return s
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.Replace(s, " ", "-", -1)
	s = strings.Replace(s, "ñ", "n", -1)
	bb := make([]byte, 0, len(s))
	for _, b := range []byte(s) {
		if (b >= 'a' && b <= 'z') ||
			(b >= '1' && b <= '9') ||
			b == '-' {
			bb = append(bb, b)
		}
	}
	return string(bb)
}
