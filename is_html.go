package web

import (
	"net/http"
	"strings"
)

func IsHTML(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	return strings.Contains(accept, "text/html")
}
