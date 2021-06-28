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

	err = app.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to run server application, %v", err)
	}
}
