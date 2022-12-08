package main

import (
	"fmt"
	"log"
	"net/http"
)

// Address to listen for incoming connections
const listenOn string = ":8080"

// Handler function to responds to an HTTP request
func HandleRequest(rw http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, req.URL.Path)

	fmt.Fprintf(rw, "Hi there!")
}

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)

	// Register handler function
	http.HandleFunc("/", HandleRequest)

	log.Printf("Listen on '%v'", listenOn)
	log.Fatal(http.ListenAndServe(listenOn, nil))
}
