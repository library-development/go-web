package web

import "net/http"

func ServeNotFound(w http.ResponseWriter, r *http.Request) {
	if IsHTML(r) {
		http.NotFound(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"error\":\"not found\"}"))
	}
}
