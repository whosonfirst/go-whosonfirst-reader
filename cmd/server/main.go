package main

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-reader/application/server"
	"log"
)

func main() {

	ctx := context.Background()

	app, err := server.NewServerApplication(ctx)

	if err != nil {
		log.Fatalf("Failed to create new server application, %v", err)
	}

	fs, err := app.DefaultFlagSet(ctx)

	if err != nil {
		log.Fatalf("Failed to create default flagset for server application, %v", err)
	}

	err = app.RunWithFlagSet(ctx, fs)

	if err != nil {
		log.Fatalf("Failed to run server application, %v", err)
	}
}
