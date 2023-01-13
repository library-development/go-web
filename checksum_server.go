package web

import (
	"net/http"
	"os"
	"path/filepath"
)

// ChecksumServer is a server for checksums.
type ChecksumServer struct {
	Dir string
}

// ServeHTTP serves HTTP requests.
func (s *ChecksumServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		path := filepath.Join(s.Dir, r.URL.Path)
		fi, err := os.Stat(path)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if fi.IsDir() {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		f, err := os.Open(path)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer f.Close()
		b := make([]byte, fi.Size())
		_, err = f.Read(b)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(b)
	case http.MethodPut:
		cksum := r.Header.Get("SHA256")
		b := make([]byte, r.ContentLength)
		_, err := r.Body.Read(b)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if cksum != checksum(b) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("checksum mismatch"))
			return
		}
		path := filepath.Join(s.Dir, r.URL.Path)
		err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		err = os.WriteFile(path, b, os.ModePerm)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
}
