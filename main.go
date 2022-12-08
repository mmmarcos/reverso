package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

// Configuration options
// TODO: Read from file or flags
const (
	listenOn   string = ":8080"          // Address to listen for incoming connections
	originHost string = "localhost:8081" // Origin server to forward requests
)

// Reverso is an HTTP handler behaving as a reverse proxy.
//
// Reverso forwards incoming requests to a target server and
// sends the response back to the client.
type Reverso struct {
	// Origin server URL to forward requests.
	originURL url.URL
}

// Handler function to responds to an HTTP request.
func (r *Reverso) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, req.URL.Path)

	// Check if origin URL is valid
	if r.originURL.Host == "" {
		r.fail(rw, &InternalError{"Origin URL is empty"}) // Fatal?
		return
	}

	// Modify request to forward to the origin server
	req.URL.Scheme = r.originURL.Scheme
	req.URL.Host = r.originURL.Host
	req.RequestURI = "" // Should be empty for client requests (see src/net/http/client.go:217)

	// Send request to the origin server
	log.Printf("Forwarding request to: '%v'", req.URL.String())
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		r.fail(rw, &InternalError{err.Error()})
		return
	}

	// Copy response from origin
	r.copyResponse(res, rw)
}

func (r *Reverso) fail(rw http.ResponseWriter, err error) {
	log.Println(err)
	rw.WriteHeader(http.StatusInternalServerError)
}

func (r *Reverso) copyResponse(res *http.Response, rw http.ResponseWriter) {
	// Copy headers
	for key, values := range res.Header {
		for _, value := range values {
			rw.Header().Add(key, value)
		}
	}
	rw.WriteHeader(res.StatusCode)

	// Write response back
	io.Copy(rw, res.Body)
}

type InternalError struct {
	msg string
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("Error: %v", e.msg)
}

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)

	// Register proxy
	http.Handle("/", &Reverso{originURL: url.URL{Scheme: "http", Host: originHost}})

	log.Printf("Listen on: '%v'", listenOn)
	log.Fatal(http.ListenAndServe(listenOn, nil))
}
