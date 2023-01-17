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
	expectedBody := "Hello from the other side"
	expectedHeaders := map[string]string{
		"X-Test-Header-1": "X-Test-Header-Value-1",
		"X-Test-Header-2": "X-Test-Header-Value-2",
	}
	expectedTrailers := map[string]string{
		"X-Test-Trailer-1": "X-Test-Trailer-Value-1",
		"X-Test-Trailer-2": "X-Test-Trailer-Value-2",
	}

	// Setup
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	svr := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		for key, value := range expectedHeaders {
			rw.Header().Set(key, value)
		}
		for key := range expectedTrailers {
			rw.Header().Add("Trailer", key)
		}

		fmt.Fprint(rw, expectedBody)

		for key, value := range expectedTrailers {
			rw.Header().Add(key, value)
		}
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

	for key, value := range expectedHeaders {
		if res.Header.Get(key) != value {
			t.Errorf("Expected header '%v: %v', got '%v'", key, value, res.Header.Get(key))
		}
	}

	if string(body) != expectedBody {
		t.Errorf("Expected body '%v', got '%v'", expectedBody, string(body))
	}

	for key, value := range expectedTrailers {
		if res.Trailer.Get(key) != value {
			t.Errorf("Expected trailer '%v: %v', got '%v'", key, value, res.Trailer.Get(key))
		}
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
