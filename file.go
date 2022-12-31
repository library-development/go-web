package web

type File struct {
	Metadata Metadata `json:"metadata"`
	Data     []byte   `json:"data"`
}

func (f *File) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if accept(r, "text/html") {
			w.Header().Set("Content-Type", "text/html")
			f.WriteHTML(w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(f)
	}
}

func (f *File) WriteHTML(w io.Writer) error {
	return f.Metadata.WriteHTML(w)
}
