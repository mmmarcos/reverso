package main

import (
	"bufio"
	"bytes"
	"log"
	"net/http"
	"sync"
	"time"
)

// Simple in-memory cache middleware.
//
// Response are indexed by the request URL path and are stored
// only if they contain the "Expires" header.
type CacheMiddleware struct {
	mutex sync.RWMutex
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

	cached, ok := c.get(req.URL.Path)

	// Response is not on cache
	if !ok {
		setCacheStatus(rw, "MISS")
		return
	}
	// Response in cache expired
	if cached.expires <= time.Now().Unix() {
		setCacheStatus(rw, "MISS")
		c.delete(req.URL.Path)
		return
	}

	// Read cached response
	res, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(cached.responseData)), req)
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

	c.update(req.URL.Path, &cachedResponse{responseData: b.Bytes(), expires: t.Unix()})
}

func (c *CacheMiddleware) get(URLPath string) (entry *cachedResponse, ok bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	entry, ok = c.cache[URLPath]
	return
}

func (c *CacheMiddleware) update(URLPath string, response *cachedResponse) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache[URLPath] = response
}

func (c *CacheMiddleware) delete(URLPath string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.cache, URLPath)
}

func setCacheStatus(rw http.ResponseWriter, status string) {
	log.Printf("Cache %v", status)
	rw.Header().Add("X-Cache-Status", status)
}
