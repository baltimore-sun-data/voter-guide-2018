package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/KenjiTakahashi/tu/titlecase"
	humanize "github.com/dustin/go-humanize"
)

const (
	results18url         = "https://elections.maryland.gov/elections/results_data/GP18/Results.js"
	precinctResults18url = "https://elections.maryland.gov/elections/results_data/GP18/PrecinctResults.js"
	metadata18url        = "https://elections.maryland.gov/elections/results_data/GP18/MetaData.js"
)

func init() {
	http.DefaultClient.Timeout = 60 * time.Second
}

func main() {
	c := FromArgs(os.Args[1:])
	if err := c.Exec(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

type Config struct {
	Local            bool
	DevServer        bool
	CreateResults    bool
	MetadataLocation string
	ResultsLocation  string
	OutputDir        string
	TemplateGlob     string
	Region           string
	Bucket           string
	Path             string
	StorySlug        string
	PollInterval     time.Duration
	CacheTime        time.Duration
	DevPort          int
	NumWorkers       int
}

func FromArgs(args []string) *Config {
	conf := &Config{}
	fl := flag.NewFlagSet("robocopy", flag.ExitOnError)
	fl.BoolVar(&conf.Local, "local", false, "just save files locally")
	fl.BoolVar(&conf.DevServer, "dev-server", false, "start a local development server")
	fl.BoolVar(&conf.CreateResults, "results", false, "create results metadata file")
	fl.StringVar(&conf.MetadataLocation, "metadata-src", metadata18url, "url or filename for metadata")
	fl.StringVar(&conf.ResultsLocation, "results-src", results18url, "url or filename for results")
	fl.StringVar(&conf.OutputDir, "output-dir", "dist/results", "directory to save into")
	fl.StringVar(&conf.TemplateGlob, "template-glob", "layouts-robocopy/*.html", "pattern to look for templates with")
	fl.StringVar(&conf.Region, "region", "us-east-1", "Amazon region for S3")
	fl.StringVar(&conf.Bucket, "bucket", "elections2018-news-baltimoresun-com", "Amazon S3 bucket")
	fl.StringVar(&conf.Path, "path", "/results/", "Amazon S3 destination path")
	fl.StringVar(&conf.StorySlug, "story-slug", "bs-2018-elections-primary-story", "story to update")
	fl.DurationVar(&conf.PollInterval, "poll-interval", 30*time.Second, "time between refreshing S3")
	fl.DurationVar(&conf.CacheTime, "cache-time", 10*time.Second, "cache time header")
	fl.IntVar(&conf.DevPort, "dev-port", 9191, "port for dev server")
	fl.IntVar(&conf.NumWorkers, "workers", 5, "number of upload workers")
	fl.Usage = func() {
		fmt.Fprintf(os.Stderr,
			`robocopy

Usage of robocopy:

`,
		)
		fl.PrintDefaults()
	}
	_ = fl.Parse(args)

	return conf
}

func (c *Config) Exec() error {
	if c.Local {
		return c.LocalExec()
	}
	if c.DevServer {
		return c.Serve()
	}
	return c.RemoteExec()
}

func (c *Config) LocalExec() error {
	t, err := c.template()
	if err != nil {
		return err
	}
	m, err := MetadataFrom(c.MetadataLocation)
	if err != nil {
		return err
	}

	if c.CreateResults {
		err = c.createJSON("results.json", m)
		if err != nil {
			return fmt.Errorf("could not create results file: %v", err)
		}
		return nil
	}

	r, err := ResultsContainerFrom(c.ResultsLocation)
	if err != nil {
		return err
	}

	cr := MapContestResults(m, r)
	for cid, rp := range cr {
		filename := filepath.Join(c.OutputDir, fmt.Sprintf("contests/%d.html", cid))
		err = c.createFile(t, "contest.html", filename, rp)
		if err != nil {
			return fmt.Errorf("could not create contest results file: %v", err)
		}
	}

	dr := MapDistrictResults(m, cr)
	for did, dp := range dr {
		filename := filepath.Join(c.OutputDir, fmt.Sprintf("districts/%d.html", did))
		err = c.createFile(t, "district.html", filename, dp)
		if err != nil {
			return fmt.Errorf("could not create district result file: %v", err)
		}
	}

	return nil
}

func (c *Config) createJSON(filename string, data interface{}) (err error) {
	os.MkdirAll(c.OutputDir, os.ModePerm)
	f, err := os.Create(filepath.Join(c.OutputDir, filename))
	if err != nil {
		return fmt.Errorf("could not create JSON file %s/%s: %v", c.OutputDir, filename, err)
	}
	defer deferClose(&err, f.Close)

	enc := json.NewEncoder(f)
	return enc.Encode(data)
}

func (c *Config) createFile(t *template.Template, tplname, path string, data interface{}) (err error) {
	dir := filepath.Dir(path)
	os.MkdirAll(dir, os.ModePerm)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create template file %s: %v", path, err)
	}
	defer deferClose(&err, f.Close)

	return t.ExecuteTemplate(f, tplname, data)
}

func LowerAlpha(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' {
			return r
		}
		if r >= 'A' && r <= 'Z' {
			return r - 'A' + 'a'
		}
		return -1
	}, s)
}

var funcMap = map[string]interface{}{
	"commas":     func(i int) string { return humanize.Comma(int64(i)) },
	"lowerAlpha": LowerAlpha,
	"len": func(i interface{}) int {
		v := reflect.ValueOf(i)
		if k := v.Kind(); k != reflect.Array &&
			k != reflect.Chan &&
			k != reflect.Map &&
			k != reflect.Slice &&
			k != reflect.String {
			return 0
		}
		return v.Len()
	},
	"titlecase": func(s string) string { return titlecase.Convert(s, nil, nil) },
	"truncateAt": func(from, s string) string {
		n := strings.Index(s, from)
		if n < 0 {
			n = len(s)
		}
		return s[:n]
	},
	"initials": func(s string) string {
		ss := strings.Fields(s)
		b := make([]byte, len(ss))
		for i, word := range ss {
			b[i] = word[0]
		}
		return string(b)
	},
}

func (c *Config) template() (*template.Template, error) {
	log.Print("getting template")
	t, err := template.New("").Funcs(funcMap).ParseGlob(c.TemplateGlob)
	if err != nil {
		return nil, fmt.Errorf("could not load templates from %s: %v", c.TemplateGlob, err)
	}
	return t, err
}
