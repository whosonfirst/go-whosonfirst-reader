package server

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-http-server"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-reader/application"
	"github.com/whosonfirst/go-whosonfirst-reader/www"
	"log"
	"net/http"
)

var server_uri string
var reader_uri string
var handler_type string

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

	handler = www.ParseURIHandler(handler)

	mux := http.NewServeMux()
	mux.Handle("/", handler)

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
