package web

import "net/http"

type DevServer struct {
	Apps map[string]http.Handler
}

func (d *DevServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := ParsePath(r.URL.Path)
	if len(path) == 0 {
		// TODO: Serve a list of apps
		http.NotFound(w, r)
		return
	}
	app, ok := d.Apps[r.Host]
	if !ok {
		http.NotFound(w, r)
		return
	}
	http.StripPrefix("/"+path[0], app).ServeHTTP(w, r)
}
