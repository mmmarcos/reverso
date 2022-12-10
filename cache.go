package main

import (
	"bufio"
	"bytes"
	"log"
	"net/http"
	"time"
)

// Basic in-memory cache for response data (currently only a named type for map)
type CacheMiddleware struct {
	cache map[string]*cachedResponse
}

type cachedResponse struct {
	responseData []byte
	expires      int64
}

// Creates an empty in-memory cache middleware
func NewCacheMiddleware() *CacheMiddleware {
	var c CacheMiddleware
	c.cache = make(map[string]*cachedResponse)
	return &c
}

// Process a request before it reaches origin server
func (c *CacheMiddleware) ProcessRequest(rw http.ResponseWriter, req *http.Request) {
	log.Printf("Processing request for: '%v'", req.URL.Path)

	entry, ok := c.cache[req.URL.Path]

	// Not on cache or expired response
	if !ok || entry.expires <= time.Now().Unix() {
		log.Printf("Cache MISS")
		rw.Header().Add("X-Cache-Status", "MISS")

		delete(c.cache, req.URL.Path) // not necessary?
		return
	}

	// Read cached (and not expired) response
	res, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(entry.responseData)), req)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Cache HIT")
	rw.Header().Add("X-Cache-Status", "HIT")
	WriteResponse(rw, res)
}

// Process a response from origin server
func (c *CacheMiddleware) ProcessResponse(res *http.Response, req *http.Request) {
	log.Printf("Processing response for: '%v'", req.URL.Path)

	// Cache response based on "Expires" header
	expiresHeader := res.Header.Get("Expires")
	if expiresHeader == "" {
		log.Println("Response does not contains Expires header, would not cache")
		return
	}

	t, err := time.Parse(http.TimeFormat, expiresHeader)
	if err != nil {
		log.Printf("Error parsing Expires header: %v", err)
		return
	}

	b := DumpResponse(res)
	c.cache[req.URL.Path] = &cachedResponse{responseData: b.Bytes(), expires: t.Unix()}
}
