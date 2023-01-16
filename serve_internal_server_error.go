package web

import "net/http"

func ServeInternalServerError(w http.ResponseWriter, r *http.Request) {
	if IsHTML(r) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"error\":\"internal server error\"}"))
	}
}
