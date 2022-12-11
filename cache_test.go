package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Tests requests resulting in HIT/MISS
func TestCacheMiddleware(t *testing.T) {
	// Setup test
	c := NewCacheMiddleware()

	// Process request for '/' on empty Cache -> Cache MISS
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c.ProcessRequest(rec, req)
	res := rec.Result()
	expected := "MISS"
	if cacheStatus := res.Header.Get("X-Cache-Status"); cacheStatus != expected {
		t.Errorf("Expected '%v', got '%v'", expected, cacheStatus)
	}

	// Process response for '/', set to expire in 1h
	freshRes := makeDummyResponse()
	freshRes.Header.Set("Expires", time.Now().In(time.FixedZone("GMT", 0)).Add(time.Hour).Format(http.TimeFormat))
	c.ProcessResponse(freshRes, req)

	// Process new request for '/' -> Cache HIT
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c.ProcessRequest(rec, req)
	res = rec.Result()
	expected = "HIT"
	if cacheStatus := res.Header.Get("X-Cache-Status"); cacheStatus != expected {
		t.Errorf("Expected '%v', got '%v'", expected, cacheStatus)
	}

	// Process response for '/', set expired
	staleRes := makeDummyResponse()
	staleRes.Header.Set("Expires", time.Now().In(time.FixedZone("GMT", 0)).Add(-time.Hour).Format(http.TimeFormat))
	c.ProcessResponse(staleRes, req)

	// Process new request for '/' -> Cache MISS
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c.ProcessRequest(rec, req)
	res = rec.Result()
	expected = "MISS"
	if cacheStatus := res.Header.Get("X-Cache-Status"); cacheStatus != expected {
		t.Errorf("Expected '%v', got '%v'", expected, cacheStatus)
	}
}

func makeDummyResponse() *http.Response {
	return &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewBufferString("Hello World")),
	}
}
