package web

import (
	"io"
	"net/http"
)

func BuildRequest(r *http.Request) (*Request, error) {
	timestamp := time.Now().UnixNano()
	body, _ := io.ReadAll(r.Body)
	return &Request{
		Timestamp: timestamp,
		FromIP:    r.RemoteAddr,
		Method:    r.Method,
		Host:      r.Host,
		Path:      r.URL.Path,
		Query:     r.URL.Query(),
		Headers:   r.Header,
		Body:      body,
	}
}
