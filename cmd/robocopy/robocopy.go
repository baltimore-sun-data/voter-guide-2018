package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
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
	MetadataLocation string
}

func FromArgs(args []string) *Config {
	conf := &Config{}
	fl := flag.NewFlagSet("robocopy", flag.ExitOnError)
	fl.BoolVar(&conf.Local, "local", false, "just save files local")
	fl.StringVar(&conf.MetadataLocation, "metadata-src", metadata18url, "url or filename for metadata")
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

	err = templates.ExecuteTemplate(os.Stdout, "metadata", m)
	return err
}
