package web

import "net/http"

func isPOST(r *http.Request) bool {
	return r.Method == http.MethodPost
}
