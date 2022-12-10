package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// Reverso is an HTTP handler behaving as a reverse proxy.
//
// Reverso forwards incoming requests to a target server and
// sends the response back to the client.
type Reverso struct {
	// Origin server URL to forward requests.
	originURL url.URL

	// In-memory cache middleware
	cache CacheMiddleware
}

// Handler function to responds to an HTTP request.
func (r *Reverso) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, req.URL.Path)

	r.cache.ProcessRequest(rw, req)

	if rw.Header().Get("X-Cache-Status") == "MISS" {

		// Fetch request from origin server
		res, err := r.fetchFromOrigin(req)
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		b := DumpResponse(res)

		r.cache.ProcessResponse(ReadResponse(b.Bytes(), req), req)

		WriteResponse(rw, ReadResponse(b.Bytes(), req))
	}
}

func (r *Reverso) fetchFromOrigin(req *http.Request) (*http.Response, error) {
	// Modify request to forward to origin server
	req.URL.Scheme = r.originURL.Scheme
	req.URL.Host = r.originURL.Host
	req.RequestURI = "" // Should be empty for client requests (see src/net/http/client.go:217)

	log.Printf("Forwarding request to: '%v'", req.URL.String())

	// Send request to the origin server
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, &internalError{err.Error()}
	}

	return res, nil
}

type internalError struct {
	msg string
}

func (e *internalError) Error() string {
	return fmt.Sprintf("Error: %v", e.msg)
}
