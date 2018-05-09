package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
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
	OutputDir        string
}

func FromArgs(args []string) *Config {
	conf := &Config{}
	fl := flag.NewFlagSet("robocopy", flag.ExitOnError)
	fl.BoolVar(&conf.Local, "local", false, "just save files local")
	fl.BoolVar(&conf.CreateResults, "results", false, "create results metadata file")
	fl.StringVar(&conf.MetadataLocation, "metadata-src", metadata18url, "url or filename for metadata")
	fl.StringVar(&conf.OutputDir, "output-dir", "static/results/", "directory to save into")
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
	if !c.Local {
		return fmt.Errorf("not implemented")
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

func (c *Config) createFile(name string, data interface{}) (err error) {
	os.MkdirAll(c.OutputDir, os.ModePerm)
	f, err := os.Create(filepath.Join(c.OutputDir, name) + ".html")
	if err != nil {
		return fmt.Errorf("could not create template file %s/%s.html: %v", c.OutputDir, name, err)
	}
	defer deferClose(&err, f.Close)

	return templates.ExecuteTemplate(f, name, data)
}
