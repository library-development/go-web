package web

import (
	"net/http"
)

// FileServer is a data structure that implements the http.Handler interface.
// It serves files and directories on the open web.
// Clients can GET and PUT files and directories.
type FileServer struct {
	AuthServerAddr string
}

// ServeHTTP serves HTTP requests for files.
// GET /path/to/file returns the file.
// PUT /path/to/file creates or updates the file.
// GET /path/to/dir returns the directory as a special type of file.
func (s *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var session *Session
	session.UserID = r.Header.Get("X-User-ID")
	session.Token = r.Header.Get("X-Token")
	err := s.authClient().ValidateSession(session)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	// path := ParsePath(r.URL.Path)
	// meta, err := s.metadata(path)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }
	// switch r.Method {
	// case http.MethodGet:
	// 	if meta.Public {
	// 		json.NewEncoder(w).Encode(meta)
	// 		return
	// 	}
	// 	// if meta.Owner {
	// 	// 	json.NewEncoder(w).Encode(meta)
	// 	// 	return
	// 	// }
	// case http.MethodPut:
	// 	// TODO
	// }
}

func (s *FileServer) authClient() *AuthClient {
	return &AuthClient{AuthServerAddr: s.AuthServerAddr}
}

// metadata returns the metadata for the given path.
// func (s *FileServer) metadata(path Path) (*Metadata, error) {
// 	resp, err := http.Get(s.MetadataServerAddr + path.String())
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	var m Metadata
// 	err = json.NewDecoder(resp.Body).Decode(&m)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &m, nil
// }
