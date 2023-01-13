package web

import (
	"encoding/json"
	"net/http"
	"time"
)

// MetadataServer is a server for global metadata.
type MetadataServer map[string]*Metadata

// Owners returns the owners of the given path.
// Deprecated.
func (s MetadataServer) Owners(path Path) map[string]bool {
	m, ok := s[path.String()]
	if !ok {
		parts := path.Parts()
		if len(parts) == 0 {
			return map[string]bool{}
		}
		return s.Owners(path[:len(parts)-1])
	}
	return m.Owners
}

// Owner returns the owner of the given path.
func (s MetadataServer) Owner(path Path) string {
	m, ok := s[path.String()]
	if !ok {
		parts := path.Parts()
		if len(parts) == 0 {
			return ""
		}
		return s.Owner(path[:len(parts)-1])
	}
	return m.Owner
}

// ServeHTTP serves HTTP requests for metadata.
func (s MetadataServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		d, ok := s[r.URL.Path]
		if !ok {
			path := ParsePath(r.URL.Path)
			owner := s.Owner(path)
			if owner == "" {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("not found"))
				return
			}
			now := time.Now().Unix()
			s[r.URL.Path] = &Metadata{
				Owner:     owner,
				CreatedAt: now,
				UpdatedAt: now,
			}
			json.NewEncoder(w).Encode(s[r.URL.Path])
			return
		}
		json.NewEncoder(w).Encode(d)
	case http.MethodPut:
		var m Metadata
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		s[r.URL.Path] = &m
	}
}
