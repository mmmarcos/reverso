package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func main() {
	const listenOn = ":8081"
	log.SetPrefix("[origin] ")
	log.Printf("Listen on: '%v'", listenOn)
	log.Fatal(http.ListenAndServe(listenOn, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Printf("Request received at: %s\n", time.Now())

		// Add "Expires" header for queries containing "s=<seconds>"
		if m, err := url.ParseQuery(req.URL.RawQuery); err == nil {
			if seconds, err := strconv.Atoi(m.Get("s")); err == nil {
				now := time.Now().In(time.FixedZone("GMT", 0))
				expires := now.Add(time.Duration(seconds) * time.Second)
				rw.Header().Set("Expires", expires.Format(http.TimeFormat))
				log.Printf("Added Header Expires: %v", rw.Header().Get("Expires"))
			}
		}
		rw.Header().Set("X-Custom-Header", "Origin")

		_, _ = fmt.Fprintf(rw, "Response from origin for path '%v'\n", req.URL.Path)
	})))
}
