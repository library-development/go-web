package web

import "net/http"

func ServeUnauthorized(w http.ResponseWriter, r *http.Request) {
	if IsHTML(r) {
		http.Error(w, "not authorized", http.StatusUnauthorized)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("{\"error\":\"not authorized\"}"))
	}
}
