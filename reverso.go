package main

import (
	"log"
	"net/http"
	"net/url"
)

// Reverso is an HTTP handler behaving as a reverse proxy.
//
// Incoming requests are forwarded to the host specified in originURL.
//
// Responses containing the "Expires" header are stored in an in-memory cache
// and served from there on further requests, as long as they do not expire
type Reverso struct {
	// Origin server URL to forward requests.
	originURL url.URL

	// In-memory cache middleware to store response data.
	cache CacheMiddleware
}

// Handler function to responds to an HTTP request.
func (r *Reverso) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, req.URL.Path)

	// Fetch from cache if available
	r.cache.ProcessRequest(rw, req) // writes X-Cache-Status header

	if rw.Header().Get("X-Cache-Status") == "MISS" {
		// Modify request host and scheme to point to origin server
		req.URL.Host = r.originURL.Host
		req.URL.Scheme = r.originURL.Scheme
		req.RequestURI = "" // Must be empty for client requests (see field docs in https://pkg.go.dev/net/http#Request)

		// Fetch from origin server
		log.Printf("Forwarding request to: '%v'", req.URL.String())
		res, err := http.DefaultTransport.RoundTrip(req)

		if err != nil {
			log.Printf("Error: %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		b := DumpResponse(res)

		// Process response, stores in cache if contains Expires header
		r.cache.ProcessResponse(ReadResponse(b.Bytes(), req), req)

		// Write response back
		WriteResponse(rw, ReadResponse(b.Bytes(), req))
	}
}
