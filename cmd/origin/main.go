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
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.SetPrefix("[origin] ")

	const addr string = ":8081"
	log.Printf("Listen address: '%s'", addr)
	log.Fatal(http.ListenAndServe(addr, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Printf("Received request: '%s'\n", req.URL.Path)

		// Parse query string
		m, err := url.ParseQuery(req.URL.RawQuery)
		if err != nil {
			msg := fmt.Sprintf("Error parsing query string: '%s'", req.URL.RawQuery)
			log.Println(msg)
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(rw, msg)
			return
		}

		// Add "Expires" header if query string contains "expires=<N>"
		if m.Has("expires") {
			now := time.Now().In(time.FixedZone("GMT", 0))
			seconds, _ := strconv.Atoi(m.Get("expires")) // 0 seconds in case of error

			// Expires N seconds from now
			expires := now.Add(time.Duration(seconds) * time.Second).Format(http.TimeFormat)
			rw.Header().Set("Expires", expires)
			log.Printf("Expires: %v", expires)
		}

		// Send chunked data if query string contains "chunked=<N>"
		if m.Has("chunked") {
			if flusher, ok := rw.(http.Flusher); ok {
				chunks, _ := strconv.Atoi(m.Get("chunked")) // 0 chunks in case of error
				for i := 1; i <= chunks; i++ {
					msg := fmt.Sprintf("Chunk %d", i)
					fmt.Fprintln(rw, msg)
					flusher.Flush()
					log.Println(msg)
					time.Sleep(1 * time.Second)
				}
			}
		}

		// Add some custom header
		rw.Header().Set("Custom-Header", "Don't panic!")

		// Response body
		fmt.Fprintln(rw, req.URL.Path)
	})))
}
