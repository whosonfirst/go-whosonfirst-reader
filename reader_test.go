package reader

import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFeatureFromID(t *testing.T) {

	ctx := context.Background()

	cwd, err := os.Getwd()

	if err != nil {
		t.Fatal(err)
	}

	fixtures := filepath.Join(cwd, "fixtures")

	r_uri := fmt.Sprintf("local://%s", fixtures)
	r, err := reader.NewReader(ctx, r_uri)

	if err != nil {
		t.Fatal(err)
	}

	wof_id := int64(101736545)

	f, err := LoadFeatureFromID(ctx, r, wof_id)

	if err != nil {
		t.Fatal(err)
	}

	if whosonfirst.Id(f) != wof_id {
		t.Fatal("Invalid WOF ID")
	}
}
