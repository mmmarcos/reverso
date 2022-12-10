package main

import (
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

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Lmicroseconds)

	// Register proxy
	http.Handle("/", &Reverso{originURL: url.URL{Scheme: "http", Host: originHost},
		cache: *NewCacheMiddleware()})

	log.Printf("Listen on: '%v'", listenOn)
	log.Fatal(http.ListenAndServe(listenOn, nil))
}
