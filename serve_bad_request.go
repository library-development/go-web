package web

import "net/http"

func ServeBadRequest(w http.ResponseWriter, r *http.Request) {
	if IsHTML(r) {
		http.Error(w, "bad request", http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"bad request\"}"))
	}
}
