package reader

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
)

func TestLoadBytes(t *testing.T) {

	ctx := context.Background()

	cwd, err := os.Getwd()

	if err != nil {
		t.Fatalf("Failed to get current working directory, %v", err)
	}

	fixtures := filepath.Join(cwd, "fixtures")

	r_uri := fmt.Sprintf("fs://%s", fixtures)
	r, err := reader.NewReader(ctx, r_uri)

	if err != nil {
		t.Fatalf("Failed to create new reader, %v", err)
	}

	wof_id := int64(101736545)

	f, err := LoadBytes(ctx, r, wof_id)

	if err != nil {
		t.Fatalf("Failed to load feature, %v", err)
	}

	id, err := properties.Id(f)

	if err != nil {
		t.Fatalf("Failed to derive ID, %v", err)
	}

	if id != wof_id {
		t.Fatal("Invalid WOF ID")
	}
}

func TestLoadFeature(t *testing.T) {

	ctx := context.Background()

	cwd, err := os.Getwd()

	if err != nil {
		t.Fatalf("Failed to get current working directory, %v", err)
	}

	fixtures := filepath.Join(cwd, "fixtures")

	r_uri := fmt.Sprintf("fs://%s", fixtures)
	r, err := reader.NewReader(ctx, r_uri)

	if err != nil {
		t.Fatalf("Failed to create new reader, %v", err)
	}

	wof_id := int64(101736545)

	f, err := LoadFeature(ctx, r, wof_id)

	if err != nil {
		t.Fatalf("Failed to load feature, %v", err)
	}

	props := f.Properties

	v, ok := props["wof:id"]

	if !ok {
		t.Fatalf("Failed to derive ID, %v", err)
	}

	id := int64(v.(float64))

	if id != wof_id {
		t.Fatal("Invalid WOF ID")
	}
}
