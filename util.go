package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net/http"
)

// Writes response into response writer
func WriteResponse(rw http.ResponseWriter, res *http.Response) {
	// Write header
	for key, values := range res.Header {
		for _, value := range values {
			rw.Header().Add(key, value)
		}
	}

	// The preferred way to send Trailers is to predeclare in the headers
	// which trailers you will later send by setting the "Trailer" header
	// to the names of the trailer keys which will come later.
	// See Examples on https://pkg.go.dev/net/http#ResponseWriter
	for key := range res.Trailer {
		rw.Header().Add("Trailer", key)
	}

	rw.WriteHeader(res.StatusCode)

	// Write body
	io.Copy(rw, res.Body)

	// Write Trailers headers
	for key, values := range res.Trailer {
		for _, value := range values {
			rw.Header().Add(key, value)
		}
	}
}

// Dumps a http.Response into a bytes.Buffer
func DumpResponse(res *http.Response) (b bytes.Buffer) {
	if err := res.Write(&b); err != nil {
		log.Println("Error dumping response to buffer")
	}
	return
}

// Read response data into http.Response
func ReadResponse(resData []byte, req *http.Request) *http.Response {
	res, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(resData)), req)
	if err != nil {
		log.Println("Error reading response from buffer")
	}
	return res
}
