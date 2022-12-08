package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHandleRequest(t *testing.T) {
	// Mock objects
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// Test
	r := &Reverso{originURL: url.URL{Scheme: "http", Host: "localhost:8081"}}
	r.ServeHTTP(rec, req)

	const expected string = "Hi there!"

	res := rec.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if string(body) != expected {
		t.Errorf("Expected '%v', got '%v'", expected, string(body))
	}
}
