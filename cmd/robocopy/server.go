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

	http.Handle(c.Path, http.StripPrefix(c.Path, http.HandlerFunc(c.devServer)))
	srv := &http.Server{Addr: port}
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

func (c *Config) devServer(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %q", r.Method, r.URL.Path)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "max-age=0")

	if !strings.HasSuffix(r.URL.Path, ".html") {
		http.NotFound(w, r)
		return
	}
	if err := c.handleRequest(w, r); err != nil {
		status := http.StatusInternalServerError
		if s, ok := err.(interface{ Status() int }); ok {
			status = s.Status()
		}
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), status)
		return
	}
	log.Println("OK")
}

type statusError int

func (s statusError) Error() string {
	return fmt.Sprintf("%d: %s", s, http.StatusText(s.Status()))
}

func (s statusError) Status() int {
	return int(s)
}

func (c *Config) handleRequest(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(strings.TrimSuffix(r.URL.Path, ".html"))
	if err != nil {
		return statusError(http.StatusNotFound)
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
