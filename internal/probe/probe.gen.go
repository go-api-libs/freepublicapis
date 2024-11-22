package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

func main() {
	start := time.Now()
	defer func() { log.Printf("Probe took %s", time.Since(start)) }()

	r, err := record()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		// Ensure recorder is saved once we're done
		if err := r.Stop(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := probe(); err != nil {
		log.Fatal(err)
	}
}

func record() (*recorder.Recorder, error) {
	r, err := recorder.NewWithOptions(&recorder.Options{
		CassetteName:       fmt.Sprintf("internal/probe/vcr/%s", time.Now().Format(time.DateOnly)),
		Mode:               recorder.ModeReplayWithNewEpisodes,
		SkipRequestLatency: true,
		RealTransport:      http.DefaultTransport,
	})
	if err != nil {
		return nil, err
	}

	// TODO: maskAuthorization
	// r.AddHook(maskAuthorization, recorder.BeforeSaveHook)

	// replace the default transport with the recorder
	http.DefaultClient.Transport = r

	return r, nil
}
