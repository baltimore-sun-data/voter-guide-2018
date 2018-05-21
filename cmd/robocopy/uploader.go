package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type client struct {
	uploader     *s3manager.Uploader
	template     *template.Template
	cachecontrol *string
	metadata     *Metadata
}

func (c *Config) RemoteExec() error {
	log.Print("connecting to AWS")
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(c.Region),
	})
	if err != nil {
		return fmt.Errorf("bad AWS credentials: %v", err)
	}
	m, err := MetadataFrom(c.MetadataLocation)
	if err != nil {
		return err
	}

	log.Printf("connecting to S3")
	var cl = client{
		uploader: s3manager.NewUploader(s),
		cachecontrol: aws.String(
			fmt.Sprintf("public, max-age=%.0f", c.PollInterval.Seconds()),
		),
		metadata: m,
	}

	ticker := time.Tick(c.PollInterval)
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)
	for {
		cl.template, err = c.template()
		if err != nil {
			return err
		}
		if err = c.RemoteTick(cl); err != nil {
			log.Printf("had errors: %v", err)
		} else {
			log.Print("finished uploading")
		}
		select {
		case <-stopChan:
			log.Print("quitting")
			return nil
		case <-ticker:
		}
	}
	panic("unreachable")
}

func (c *Config) RemoteTick(cl client) error {
	var (
		funcCh     = make(chan func() error, c.NumWorkers)
		errCh      = make(chan error, c.NumWorkers)
		waitingFor int
	)

	for i := 0; i < c.NumWorkers; i++ {
		go func() {
			for f := range funcCh {
				errCh <- f()
			}
		}()
	}

	r, err := ResultsContainerFrom(c.ResultsLocation)
	if err != nil {
		return err
	}

	log.Print("processing and uploading results")
	cr := MapContestResults(cl.metadata, r)
	// make a queue of cids
	cids := make([]ContestID, 0, len(cr))
	for cid := range cr {
		cids = append(cids, cid)
	}

	log.Printf("received %d items", len(cids))

	// loop through and pop off cids until they're all gone
	var (
		loops  int
		hadErr error
	)
	for len(cids) > 0 || waitingFor > 0 {
		var (
			taskCh chan func() error
			task   func() error
		)
		if len(cids) > 0 {
			cid := cids[0]
			rp := cr[cid]
			filename := fmt.Sprintf("%d.html", cid)
			taskCh = funcCh
			task = func() error { return c.uploadFile(cl, filename, rp) }
		}

		select {
		case taskCh <- task:
			waitingFor++
			cids = cids[1:]
		case err := <-errCh:
			loops++
			waitingFor--
			if err != nil && hadErr == nil {
				hadErr = err
				cids = nil // Just give up for now
			}
		}
	}
	log.Printf("handled %d items", loops)
	return hadErr
}

func (c *Config) uploadFile(cl client, filename string, data interface{}) error {
	var buf = &bytes.Buffer{}
	err := cl.template.ExecuteTemplate(buf, "contest.html", data)
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	_, err = cl.uploader.Upload(&s3manager.UploadInput{
		Bucket:       aws.String(c.Bucket),
		Key:          aws.String(c.Path + filename),
		ContentType:  aws.String("text/html; charset=utf-8"),
		CacheControl: cl.cachecontrol,
		Body:         buf,
	})
	if err != nil {
		log.Printf("error uploading %s: %v", filename, err)
	}
	return err
}
