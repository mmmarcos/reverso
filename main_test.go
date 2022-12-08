package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleRequest(t *testing.T) {
	// Mock objects
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// Test
	HandleRequest(rec, req)

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
