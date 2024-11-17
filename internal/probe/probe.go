package main

import (
	"net/http"

	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

func probe() error {
	// Start our recorder
	r, err := recorder.NewWithOptions(&recorder.Options{
		CassetteName:       "fixtures/vcr",
		Mode:               recorder.ModeReplayWithNewEpisodes,
		SkipRequestLatency: true,
		RealTransport:      http.DefaultTransport,
	})
	if err != nil {
		return err
	}
	defer r.Stop() // Ensure recorder is stopped once done with it

	// TODO: maskAuthorization
	// r.AddHook(maskAuthorization, recorder.BeforeSaveHook)

	// Inject the transport
	// notion.SetTransport(r)

	return nil
}
