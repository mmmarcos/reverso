package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// Checks if client receives body and headers sent from origin
func TestHandleSimpleRequest(t *testing.T) {
	const expectedBodyStr string = "Hello from the other side"
	const customHeaderKey string = "X-Test-Header"
	const customHeaderVal string = "Custom header from origin"

	// Mock objects
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(customHeaderKey, customHeaderVal)
		fmt.Fprint(w, expectedBodyStr)
	}))
	defer svr.Close()

	// Test reverse proxy
	r := &Reverso{originURL: parseServerURL(svr.URL)}
	r.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	bodyStr := readAll(res.Body)
	if bodyStr != expectedBodyStr {
		t.Errorf("Expected '%v', got '%v'", expectedBodyStr, bodyStr)
	}

	if v := res.Header.Get(customHeaderKey); v != customHeaderVal {
		t.Errorf("Expected header value '%v', got '%v'", customHeaderVal, v)
	}
}

// Checks if client receives HTTP status code sent by origin
func TestHandleTeapotRequest(t *testing.T) {
	// Mock objects
	req := httptest.NewRequest(http.MethodGet, "/teapot", nil)
	rec := httptest.NewRecorder()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	defer svr.Close()

	// Test reverse proxy
	r := &Reverso{originURL: parseServerURL(svr.URL)}
	r.ServeHTTP(rec, req)

	res := rec.Result()

	if res.StatusCode != http.StatusTeapot {
		t.Errorf("Expected status code '%v', got '%v'", http.StatusTeapot, res.StatusCode)
	}
}

// Checks for internal server errors when using an invalid server URL
func TestInvalidOriginURL(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	r := &Reverso{} // empty origin URL
	r.ServeHTTP(rec, req)
	res := rec.Result()

	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status code '%v', got '%v'", http.StatusInternalServerError, res.StatusCode)
	}
}

// Parses the given raw URL into a url.URL
func parseServerURL(rawURL string) url.URL {
	serverURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal(err)
	}
	return *serverURL
}

// Read all from r and return content as string
func readAll(r io.Reader) string {
	body, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}
