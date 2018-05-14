package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	humanize "github.com/dustin/go-humanize"
)

const (
	results18url         = "http://elections.maryland.gov/elections/results_data/GP18/Results.js"
	precinctResults18url = "http://elections.maryland.gov/elections/results_data/GP18/PrecinctResults.js"
	metadata18url        = "http://elections.maryland.gov/elections/results_data/GP18/MetaData.js"
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
	CreateResults    bool
	MetadataLocation string
	ResultsLocation  string
	OutputDir        string
	TemplateGlob     string
	Region           string
	Bucket           string
	Path             string
	PollInterval     time.Duration
}

func FromArgs(args []string) *Config {
	conf := &Config{}
	fl := flag.NewFlagSet("robocopy", flag.ExitOnError)
	fl.BoolVar(&conf.Local, "local", false, "just save files local")
	fl.BoolVar(&conf.CreateResults, "results", false, "create results metadata file")
	fl.StringVar(&conf.MetadataLocation, "metadata-src", metadata18url, "url or filename for metadata")
	fl.StringVar(&conf.ResultsLocation, "results-src", results18url, "url or filename for results")
	fl.StringVar(&conf.OutputDir, "output-dir", "dist/results/contests", "directory to save into")
	fl.StringVar(&conf.TemplateGlob, "template-glob", "layouts-robocopy/*.html", "pattern to look for templates with")
	fl.StringVar(&conf.Region, "region", "us-east-1", "Amazon region for S3")
	fl.StringVar(&conf.Bucket, "bucket", "elections2018-news-baltimoresun-com", "Amazon S3 bucket")
	fl.StringVar(&conf.Path, "path", "/results/contests/", "Amazon S3 destination path")
	fl.DurationVar(&conf.PollInterval, "poll-interval", 30*time.Second, "time between refreshing S3")
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
	t, err := c.template()
	if err != nil {
		return err
	}
	m, err := MetadataFrom(c.MetadataLocation)
	if err != nil {
		return err
	}

	if !c.Local {
		return c.Remote(m, t)
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

func (c *Config) Remote(m *Metadata, t *template.Template) error {
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(c.Region),
	})
	if err != nil {
		return fmt.Errorf("bad AWS credentials: %v", err)
	}
	var cl = client{
		uploader: s3manager.NewUploader(s),
		cachecontrol: aws.String(
			fmt.Sprintf("public, max-age=%.0f", c.PollInterval.Seconds()),
		),
		template: t,
	}

	for range time.Tick(c.PollInterval) {
		log.Print("getting results")
		r, err := ResultsContainerFrom(c.ResultsLocation)
		if err != nil {
			log.Printf("results error: %v", err)
			continue
		}

		cr := MapContestResults(m, r)
		for cid, rp := range cr {
			filename := fmt.Sprintf("%d.html", cid)
			err = c.uploadFile(cl, filename, rp)
			if err != nil {
				return err
			}
		}
		log.Print("finished uploading")
	}
	panic("unreachable")
}

var funcMap = map[string]interface{}{
	"commas": func(i int) string { return humanize.Comma(int64(i)) },
}

func (c *Config) template() (*template.Template, error) {
	t, err := template.New("").Funcs(funcMap).ParseGlob(c.TemplateGlob)
	if err != nil {
		return nil, fmt.Errorf("could not load templates from %s: %v", c.TemplateGlob, err)
	}
	return t, err
}

type client struct {
	uploader     *s3manager.Uploader
	template     *template.Template
	cachecontrol *string
}

func (c *Config) uploadFile(cl client, filename string, data interface{}) error {
	var buf bytes.Buffer
	err := cl.template.ExecuteTemplate(&buf, "contest.html", data)
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	_, err = cl.uploader.Upload(&s3manager.UploadInput{
		Bucket:       aws.String(c.Bucket),
		Key:          aws.String(c.Path + filename),
		CacheControl: cl.cachecontrol,
		Body:         &buf,
	})
	if err != nil {
		log.Printf("error uploading %s: %v", filename, err)
	}
	return nil
}
