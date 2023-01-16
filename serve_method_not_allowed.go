package web

import "net/http"

func ServeMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	if IsHTML(r) {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("{\"error\":\"method not allowed\"}"))
	}
}
