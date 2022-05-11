package reader

import (
	"context"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"compress/bzip2"
	"fmt"
	"github.com/whosonfirst/go-ioutil"
	"strconv"
)

type FileReader struct {
	Reader
	root string
	allow_bz2 bool
}

func init() {

	ctx := context.Background()

	err := RegisterReader(ctx, "fs", NewFileReader) // Deprecated

	if err != nil {
		panic(err)
	}

}

func NewFileReader(ctx context.Context, uri string) (Reader, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	root := u.Path
	info, err := os.Stat(root)

	if err != nil {
		return nil, fmt.Errorf("Failed to stat %s, %w", root, err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("Root (%s) is not a directory", root)
	}
	
	r := &FileReader{
		root: root,
	}

	q := u.Query()

	allow_bz2 := q.Get("allow_bz2")	

	if allow_bz2 != "" {
		
		allow, err := strconv.ParseBool(allow_bz2)

		if err != nil {
			return nil, fmt.Errorf("Unable to parse '%s' parameter, %w", allow_bz2, err)
		}

		r.allow_bz2 = allow
	}
	
	return r, nil
}

func (r *FileReader) Read(ctx context.Context, path string) (io.ReadSeekCloser, error) {

	abs_path := r.ReaderURI(ctx, path)

	_, err := os.Stat(abs_path)

	if err != nil {
		return nil, fmt.Errorf("Failed to stat %s, %v", abs_path, err)
	}

	var fh io.ReadSeekCloser
	
	fh, err = os.Open(abs_path)

	if err != nil {
		return nil, fmt.Errorf("Failed to open %s, %w", abs_path, err)
	}

	if filepath.Ext(abs_path) == ".bz2" && r.allow_bz2 {

		bz_r := bzip2.NewReader(fh)
		
		rsc, err := ioutil.NewReadSeekCloser(bz_r)

		if err != nil {
			return nil, fmt.Errorf("Failed create ReadSeekCloser for bzip2 reader for %s, %w", path, err)
		}

		fh = rsc
	}

	return fh, nil
}

func (r *FileReader) ReaderURI(ctx context.Context, path string) string {
	return filepath.Join(r.root, path)
}
