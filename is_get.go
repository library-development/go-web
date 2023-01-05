package web

import "net/http"

func isGET(r *http.Request) bool {
	return r.Method == http.MethodGet
}
