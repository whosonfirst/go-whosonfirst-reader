// package server provides an HTTP server application for emitting (reading) data from a go-reader instance
package server

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-http-cache"
	"github.com/aaronland/go-http-cache/adapter/memory"
	"github.com/aaronland/go-http-server"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-reader/application"
	"github.com/whosonfirst/go-whosonfirst-reader/www"
	"log"
	"net/http"
	"time"
)

var server_uri string
var reader_uri string
var handler_type string
var enable_cache bool

type ServerApplication struct {
	application.Application
}

func NewServerApplication(ctx context.Context) (application.Application, error) {
	app := &ServerApplication{}
	return app, nil
}

func (app *ServerApplication) DefaultFlagSet(ctx context.Context) (*flag.FlagSet, error) {

	fs := flagset.NewFlagSet("reader")

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "A valid aaronland/go-http-server URI.")
	fs.StringVar(&reader_uri, "reader-uri", "null://", "A valid whosonfirst/go-reader URI.")
	fs.StringVar(&handler_type, "handler-type", "data", "Valid options are: data, redirect.")
	fs.BoolVar(&enable_cache, "enable-cache", true, "Enable the aaronland/go-http-cache middleware handler for caching lookups.")

	return fs, nil
}

func (app *ServerApplication) Run(ctx context.Context) error {

	fs, err := app.DefaultFlagSet(ctx)

	if err != nil {
		return err
	}

	return app.RunWithFlagSet(ctx, fs)
}

func (app *ServerApplication) RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVarsWithFeedback(fs, "SERVER", true)

	if err != nil {
		return fmt.Errorf("Failed to set flags from environment variables, %w", err)
	}

	r, err := reader.NewReader(ctx, reader_uri)

	if err != nil {
		return fmt.Errorf("Failed to create new reader, %v", err)
	}

	var handler http.HandlerFunc

	switch handler_type {
	case "data":

		data_handler, err := www.DataHandler(r)

		if err != nil {
			return fmt.Errorf("Failed to create new reader handler, %v", err)
		}

		handler = data_handler

	case "redirect":

		redirect_handler, err := www.RedirectHandler(r)

		if err != nil {
			return fmt.Errorf("Failed to create new redirect handler, %v", err)
		}

		handler = redirect_handler

	default:
		return fmt.Errorf("Invalid -handler-type (%s)", handler_type)
	}

	// Cache the previous handler (data, redirect) results since they might be
	// expensive especially if we are using a go-whosonfirst-findingaid instance
	// the reader.Reader itself to locate a record (to read).
	
	if enable_cache {

		memcached, err := memory.NewAdapter(
			memory.AdapterWithAlgorithm(memory.LRU),
			memory.AdapterWithCapacity(10000000),	// TO DO: make this a CLI option
		)
		
		if err != nil {
			return fmt.Errorf("Failed to create cache memory adapter, %w", err)
		}
		
		cacheClient, err := cache.NewClient(
			cache.ClientWithAdapter(memcached),
			cache.ClientWithTTL(10*time.Minute),	// TO DO: make this a CLI option
			cache.ClientWithRefreshKey("opn"),	// TO DO: make this a CLI option
		)
		
		if err != nil {
			return fmt.Errorf("Failed to create cache client, %w", err)
		}
		
		handler = cacheClient.Middleware(handler)
	}

	// This one is important: This is the first handler that will be invoked that will
	// resolve the request URI in to a valid Who's On First relative path. That path will
	// be stored in the handler's `X-WhosOnFirst-Rel-Path` response header. It is assumed
	// the header will be read and acted on by "downstream" middleware handler.
	
	handler = www.ParseURIHandler(handler)

	mux := http.NewServeMux()
	mux.Handle("/", handler)

	// For example: This is the same ParseURIHandler as above but without any "next"
	// handlers and so the default behaviour is simply to write the `X-WhosOnFirst-Rel-Path`
	// response header and print the relative path to the browser.
	
	uri_handler := www.ParseURIHandler(nil)
	
	mux.Handle("/uri/", uri_handler)

	s, err := server.NewServer(ctx, server_uri)

	if err != nil {
		return fmt.Errorf("Failed to create new server, %s", err)
	}

	log.Printf("Listening on %s\n", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		return fmt.Errorf("Failed to serve requests, %s", err)
	}

	return nil
}
