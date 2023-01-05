package web

import (
	"net/http"
	"strings"
)

func isHTML(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	return strings.Contains(accept, "text/html")
}
