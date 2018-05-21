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
	PollInterval     time.Duration
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
	fl.StringVar(&conf.OutputDir, "output-dir", "dist/results/contests", "directory to save into")
	fl.StringVar(&conf.TemplateGlob, "template-glob", "layouts-robocopy/*.html", "pattern to look for templates with")
	fl.StringVar(&conf.Region, "region", "us-east-1", "Amazon region for S3")
	fl.StringVar(&conf.Bucket, "bucket", "elections2018-news-baltimoresun-com", "Amazon S3 bucket")
	fl.StringVar(&conf.Path, "path", "/results/contests/", "Amazon S3 destination path")
	fl.DurationVar(&conf.PollInterval, "poll-interval", 30*time.Second, "time between refreshing S3")
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
		filename := fmt.Sprintf("%d.html", cid)
		err = c.createFile(t, "contest.html", filename, rp)
		if err != nil {
			return fmt.Errorf("could not create results file: %v", err)
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

func (c *Config) createFile(t *template.Template, tplname, filename string, data interface{}) (err error) {
	os.MkdirAll(c.OutputDir, os.ModePerm)
	f, err := os.Create(filepath.Join(c.OutputDir, filename))
	if err != nil {
		return fmt.Errorf("could not create template file %s/%s: %v", c.OutputDir, filename, err)
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
}

func (c *Config) template() (*template.Template, error) {
	log.Print("getting template")
	t, err := template.New("").Funcs(funcMap).ParseGlob(c.TemplateGlob)
	if err != nil {
		return nil, fmt.Errorf("could not load templates from %s: %v", c.TemplateGlob, err)
	}
	return t, err
}
