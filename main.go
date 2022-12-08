package main

import (
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
		log.Fatal("Origin URL is empty")
	}

	// Modify request to forward to the origin server
	req.URL.Scheme = r.originURL.Scheme
	req.URL.Host = r.originURL.Host
	req.RequestURI = "" // Should be empty for client requests (see src/net/http/client.go:217)

	// Send request to the origin server
	log.Printf("Forwarding request to: '%v'", req.URL.String())
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Printf("ERROR: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write response back
	rw.WriteHeader(res.StatusCode)
	io.Copy(rw, res.Body)
}

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)

	// Register proxy
	http.Handle("/", &Reverso{originURL: url.URL{Scheme: "http", Host: originHost}})

	log.Printf("Listen on: '%v'", listenOn)
	log.Fatal(http.ListenAndServe(listenOn, nil))
}
