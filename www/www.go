// package www provides common net/http handlers for reader-related HTTP requests.
package www

import (
	"github.com/whosonfirst/go-reader"
	uri "github.com/whosonfirst/go-whosonfirst-uri/http"
	"io"
	_ "log"
	"net/http"
)

// Emit an HTTP redirect header for the value of r.ReaderURI for a URI. This handler is meant to be used
// in conjunction with ParseURIHandler middleware handler. For example:
// 	handler, _ := www.RedirectHandler(r)
//	handler = www.ParseURIHandler(handler)
func RedirectHandler(r reader.Reader) (http.HandlerFunc, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		rel_path := rsp.Header().Get(uri.HEADER_RELPATH)

		if rel_path == "" {
			http.Error(rsp, "Unable to determine URI", http.StatusNotFound)
			return
		}

		abs_path := r.ReaderURI(ctx, rel_path)
		http.Redirect(rsp, req, abs_path, http.StatusFound)
		return
	}

	h := http.HandlerFunc(fn)
	return h, nil
}

// Emit the out of r.Read for a URI. This handler is meant to be used
// in conjunction with ParseURIHandler middleware handler. For example:
// 	handler, _ := www.DataHandler(r)
//	handler = www.ParseURIHandler(handler)
func DataHandler(r reader.Reader) (http.HandlerFunc, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		rel_path := rsp.Header().Get(uri.HEADER_RELPATH)

		if rel_path == "" {
			http.Error(rsp, "Unable to determine URI", http.StatusNotFound)
			return
		}

		fh, err := r.Read(ctx, rel_path)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-Type", "application/json")

		_, err = io.Copy(rsp, fh)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
