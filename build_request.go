package web

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

// Build request returns a newly created Request object.
// It leaves the request body open for further processing.
func BuildRequest(r *http.Request) (*Request, error) {
	timestamp := time.Now().UnixNano()
	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	return &Request{
		Timestamp: timestamp,
		FromIP:    r.RemoteAddr,
		Method:    r.Method,
		Host:      r.Host,
		Path:      r.URL.Path,
		Query:     r.URL.Query(),
		Headers:   r.Header,
		Body:      body,
	}, nil
}
