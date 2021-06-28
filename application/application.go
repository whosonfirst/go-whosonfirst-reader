// package application provides common logic for whosonfirst/go-reader related applications targeting Who's On First documents.
package application

import (
	"context"
	"flag"
)

type Application interface {
	DefaultFlagSet(context.Context) (*flag.FlagSet, error)
	Run(context.Context) error
	RunWithFlagSet(context.Context, *flag.FlagSet) error
}
