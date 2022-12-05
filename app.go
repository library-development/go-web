package web

import "net/http"

type App struct {
	Models []Model `json:"models"`
}

func (a *App) Serve(port, dataDir string) error {
	mux := http.NewServeMux()
	mux.Handle("/", &AppInstance{DataDir: dataDir, Source: a})
	return http.ListenAndServe(":"+port, mux)
}
