package main

import (
	"bufio"
	"bytes"
	"log"
	"net/http"
	"time"
)

// Simple in-memory cache middleware.
//
// Response are indexed by the request URL path and are stored
// only if they contain the "Expires" header.
type CacheMiddleware struct {
	cache map[string]*cachedResponse
}

// A cached response with its given Expires
type cachedResponse struct {
	responseData []byte
	expires      int64
}

// Creates an empty cache middleware
func NewCacheMiddleware() *CacheMiddleware {
	var c CacheMiddleware
	c.cache = make(map[string]*cachedResponse)
	return &c
}

// Process a request before it reaches origin server
func (c *CacheMiddleware) ProcessRequest(rw http.ResponseWriter, req *http.Request) {
	log.Printf("Processing request for: '%v'", req.URL.Path)

	entry, ok := c.cache[req.URL.Path]

	// Response is not on cache or already expired
	if !ok || entry.expires <= time.Now().Unix() {
		setCacheStatus(rw, "MISS")
		delete(c.cache, req.URL.Path) // not necessary?
		return
	}

	// Read cached response
	res, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(entry.responseData)), req)
	if err != nil {
		log.Printf("Error reading response stored in cache: %v", err)
		setCacheStatus(rw, "MISS")
		return
	}

	setCacheStatus(rw, "HIT")
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

func setCacheStatus(rw http.ResponseWriter, status string) {
	log.Printf("Cache %v", status)
	rw.Header().Add("X-Cache-Status", status)
}
