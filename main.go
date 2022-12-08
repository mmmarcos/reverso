package main

import (
	"fmt"
	"log"
	"net/http"
)

// Address to listen for incoming connections
const listenOn string = ":8080"

// Reverso is an HTTP handler behaving as a reverse proxy.
//
// Reverso forwards incoming requests to a target server and
// sends the response back to the client.
type Reverso struct {
}

// Handler function to responds to an HTTP request.
func (r *Reverso) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, req.URL.Path)

	fmt.Fprintf(rw, "Hi there!")
}

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)

	// Register proxy
	http.Handle("/", &Reverso{})

	log.Printf("Listen on '%v'", listenOn)
	log.Fatal(http.ListenAndServe(listenOn, nil))
}
