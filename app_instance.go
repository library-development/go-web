package web

import "net/http"

type AppInstance struct {
	DataDir string `json:"data_dir"`
	Source  *App   `json:"source"`
}

func (a *AppInstance) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
