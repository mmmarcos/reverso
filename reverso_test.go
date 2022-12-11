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

// Checks if client receives response sent by origin
func TestHandleSimpleRequest(t *testing.T) {
	const expectedBodyStr string = "Hello from the other side"
	const customHeaderKey string = "X-Test-Header"
	const customHeaderVal string = "Custom header from origin"

	// Setup
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(customHeaderKey, customHeaderVal)
		fmt.Fprint(w, expectedBodyStr)
	}))
	defer svr.Close()
	r := &Reverso{originURL: parseServerURL(svr.URL), cache: *NewCacheMiddleware()}

	// Test reverse proxy handler
	r.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	if string(body) != expectedBodyStr {
		t.Errorf("Expected '%v', got '%v'", expectedBodyStr, string(body))
	}

	if v := res.Header.Get(customHeaderKey); v != customHeaderVal {
		t.Errorf("Expected header '%v:%v', got '%v'", customHeaderKey, customHeaderVal, v)
	}
}

// Checks if client receives HTTP status code sent by origin
func TestHandleTeapotRequest(t *testing.T) {
	// Setup
	req := httptest.NewRequest(http.MethodGet, "/teapot", nil)
	rec := httptest.NewRecorder()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	defer svr.Close()

	// Test reverse proxy handler
	r := &Reverso{originURL: parseServerURL(svr.URL), cache: *NewCacheMiddleware()}
	r.ServeHTTP(rec, req)

	res := rec.Result()

	if res.StatusCode != http.StatusTeapot {
		t.Errorf("Expected status code '%v', got '%v'", http.StatusTeapot, res.StatusCode)
	}
}

func TestInvalidOriginURL(t *testing.T) {
	// Setup
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	r := &Reverso{} // empty origin URL

	// Test reverse proxy handler
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
