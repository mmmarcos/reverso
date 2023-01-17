package main

import (
	"log"
	"net/http"
	"net/url"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Lmicroseconds)

	http.Handle("/", &Reverso{originURL: url.URL{Scheme: "http", Host: "localhost:8081"}, cache: *NewCacheMiddleware()})

	const addr string = ":8080"
	log.Printf("Listen address: '%s'", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
