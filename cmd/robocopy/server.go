package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

func (c *Config) Serve() error {
	port := fmt.Sprintf(":%d", c.DevPort)

	// subscribe to SIGINT signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	srv := &http.Server{Addr: port, Handler: c.Routes()}
	go func() {
		<-stopChan // wait for system signal
		log.Println("Shutting down server...")

		// shut down gracefully, but wait no longer than 5 seconds before halting
		ctx, c := context.WithTimeout(context.Background(), 5*time.Second)
		defer c()
		srv.Shutdown(ctx)
	}()

	log.Printf("Serving http://localhost:%d%s", c.DevPort, c.Path)
	return srv.ListenAndServe()
}

func (c *Config) Routes() http.Handler {
	mux := http.NewServeMux()
	for _, br := range BarkerResults(nil) {
		mux.Handle(c.Path+br.Slug, http.StripPrefix(c.Path,
			c.stdMiddleware(c.handleBarker)))
	}
	mux.Handle(c.Path+c.StorySlug, http.StripPrefix(c.Path+c.StorySlug,
		c.stdMiddleware(c.handleStory)))
	mux.Handle(c.Path+"contests/", http.StripPrefix(c.Path+"contests/",
		c.stdMiddleware(c.handleContests)))
	mux.Handle(c.Path+"districts/", http.StripPrefix(c.Path+"districts/",
		c.stdMiddleware(c.handleDistricts)))
	return mux
}

type erroringHandler = func(w http.ResponseWriter, r *http.Request) error

func (c *Config) stdMiddleware(h erroringHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %q", r.Method, r.URL.Path)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Cache-Control", "max-age=0")

		if err := h(w, r); err != nil {
			status := http.StatusInternalServerError
			if s, ok := err.(interface{ Status() int }); ok {
				status = s.Status()
			}
			log.Printf("error: %v", err)
			http.Error(w, err.Error(), status)
			return
		}
		log.Println("OK")
	})
}

type statusError int

func (s statusError) Error() string {
	return fmt.Sprintf("%d: %s", s, http.StatusText(s.Status()))
}

func (s statusError) Status() int {
	return int(s)
}

func extractID(r *http.Request) (int, error) {
	if !strings.HasSuffix(r.URL.Path, ".html") {
		return 0, statusError(http.StatusNotFound)
	}
	id, err := strconv.Atoi(strings.TrimSuffix(r.URL.Path, ".html"))
	if err != nil {
		return 0, statusError(http.StatusNotFound)
	}
	return id, nil
}

func (c *Config) handleContests(w http.ResponseWriter, r *http.Request) error {
	id, err := extractID(r)
	if err != nil {
		return err
	}
	cid := ContestID(id)

	t, err := c.template()
	if err != nil {
		return err
	}
	m, err := MetadataFrom(c.MetadataLocation)
	if err != nil {
		return err
	}

	rc, err := ResultsContainerFrom(c.ResultsLocation)
	if err != nil {
		return err
	}

	cr := MapContestResults(m, rc)
	data, ok := cr[cid]
	if !ok {
		return statusError(http.StatusNotFound)
	}

	err = t.ExecuteTemplate(w, "contest.html", data)
	if err != nil {
		// too late to return 500
		log.Printf("error executing template: %v", err)
	}
	return nil
}

func (c *Config) handleDistricts(w http.ResponseWriter, r *http.Request) error {
	id, err := extractID(r)
	if err != nil {
		return err
	}
	did := DistrictID(id)

	t, err := c.template()
	if err != nil {
		return err
	}
	m, err := MetadataFrom(c.MetadataLocation)
	if err != nil {
		return err
	}

	rc, err := ResultsContainerFrom(c.ResultsLocation)
	if err != nil {
		return err
	}

	cr := MapContestResults(m, rc)
	dr := MapDistrictResults(m, cr)
	data, ok := dr[did]
	if !ok {
		return statusError(http.StatusNotFound)
	}

	err = t.ExecuteTemplate(w, "district.html", data)
	if err != nil {
		// too late to return 500
		log.Printf("error executing template: %v", err)
	}
	return nil
}

func (c *Config) handleBarker(w http.ResponseWriter, r *http.Request) error {
	m, err := MetadataFrom(c.MetadataLocation)
	if err != nil {
		return err
	}

	rc, err := ResultsContainerFrom(c.ResultsLocation)
	if err != nil {
		return err
	}
	cr := MapContestResults(m, rc)
	brs := BarkerResults(cr)
	for _, br := range brs {
		if r.URL.Path == br.Slug {
			return c.handleP2P("barker.html", br, w, r)
		}
	}
	return statusError(http.StatusNotFound)
}

func (c *Config) handleStory(w http.ResponseWriter, r *http.Request) error {
	m, err := MetadataFrom(c.MetadataLocation)
	if err != nil {
		return err
	}

	rc, err := ResultsContainerFrom(c.ResultsLocation)
	if err != nil {
		return err
	}
	cr := MapContestResults(m, rc)

	return c.handleP2P("story.html", cr, w, r)
}

func (c *Config) handleP2P(templatename string, data interface{}, w http.ResponseWriter, r *http.Request) error {

	w.Header().Set("Content-Type", "text/html")

	t, err := c.template()
	if err != nil {
		return err
	}

	err = t.ExecuteTemplate(w, templatename, data)
	if err != nil {
		// too late to return 500
		log.Printf("error executing template: %v", err)
	}
	return nil
}
