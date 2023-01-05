package web

import (
	"net/http"
	"strings"
)

func accept(r *http.Request, accept string) bool {
	return strings.Contains(r.Header.Get("Accept"), accept)
}
