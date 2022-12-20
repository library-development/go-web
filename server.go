package web

type Server struct {
	FileDir string
	LogsDir string
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := logRequest(r, h.logsDir); err != nil {
		panic(err)
	}
	host := r.Host
	dir := filepath.Join(h.dir, host)
	if r.Method == http.MethodGet {
		http.FileServer(http.Dir(dir)).ServeHTTP(w, r)
	} else {
		writeSuccess(w, r)
	}
}
