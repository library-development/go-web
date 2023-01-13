package web

import (
	"net/http"
	"strings"
)

type App struct {
	Auth *AuthDB
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/auth") {
		http.StripPrefix("/auth", a.Auth).ServeHTTP(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/docs") {
		return
	}
	if strings.HasPrefix(r.URL.Path, "/api") {
		return
	}
}
