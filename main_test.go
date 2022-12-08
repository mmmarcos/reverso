package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHandleRequest(t *testing.T) {
	const expected string = "Hello from the other side"

	// Mock objects
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	}))
	defer svr.Close()

	// Parse origin server URL
	fmt.Println(string(svr.URL))
	serverURL, err := url.Parse(svr.URL)
	if err != nil {
		t.Error(err)
	}

	// Test reverse proxy
	r := &Reverso{originURL: *serverURL}
	r.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != expected {
		t.Errorf("Expected '%v', got '%v'", expected, string(body))
	}
}
