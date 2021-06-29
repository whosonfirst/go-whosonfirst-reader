# go-http-cache

This is a high performance Golang HTTP middleware for server-side application layer caching, ideal for REST APIs.

The memory adapter minimizes GC overhead to near zero and supports some options of caching algorithms (LRU, MRU, LFU, MFU). This way, it is able to store plenty of gigabytes of responses, keeping great performance and being free of leaks.

## Important

This is a hard fork of @victorspinger's [http-cache](https://github.com/victorspringer/http-cache) package. Differences include:

* Removing the `memory/redis` adapter and the `benchmark` package so there are no external dependencies.
* Ensuring that all the response headers from the previous response are assigned to the HTTP recorder.
* Update `client.Middleware` method to return `http.HandlerFunc`.
* Storing the HTTP status code in the cached `Response` and checking for and issuing HTTP redirect responses where appropriate.

## Getting Started

### Installation
`go get github.com/aaronland/go-http-cache`

### Usage
This is an example of use with the memory adapter:

```go
package main

import (
    "fmt"
    "net/http"
    "os"
    "time"
    
    "github.com/aaronland/go-http-cache"
    "github.com/aaronland/go-http-cache/adapter/memory"
)

func example(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Ok"))
}

func main() {
    memcached, err := memory.NewAdapter(
        memory.AdapterWithAlgorithm(memory.LRU),
        memory.AdapterWithCapacity(10000000),
    )
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    cacheClient, err := cache.NewClient(
        cache.ClientWithAdapter(memcached),
        cache.ClientWithTTL(10 * time.Minute),
        cache.ClientWithRefreshKey("opn"),
    )
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    handler := http.HandlerFunc(example)

    http.Handle("/", cacheClient.Middleware(handler))
    http.ListenAndServe(":8080", nil)
}
```


## Godoc Reference
- [http-cache](https://godoc.org/github.com/aaronland/go-http-cache)
- [Memory adapter](https://godoc.org/github.com/aaronland/go-http-cache/adapter/memory)

## License

http-cache is released under the [MIT License](https://github.com/aaronland/go-http-cache/blob/master/LICENSE).