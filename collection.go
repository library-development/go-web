package web

import (
	"encoding/json"
	"io"
	"net/http"
)

type Collection[T http.Handler] map[string]T

func (c Collection[T]) Get(id string) (T, bool) {
	item, ok := c[id]
	return item, ok
}

func (c Collection[T]) Post(v T) string {
	id := NewID()
	c[id] = v
	return id
}

func (c Collection[T]) Put(id string, v T) {
	c[id] = v
}

func (c Collection[T]) Delete(id string) {
	delete(c, id)
}

func (c Collection[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := ParsePath(r.URL.Path)
	if len(path) == 0 {
		c.serveRoot(w, r)
		return
	}
	id := path[0]
	item, ok := c[id]
	if !ok {
		ServeNotFound(w, r)
		return
	}
	if len(path) == 1 {
		if r.Method == http.MethodDelete {
			c.Delete(id)
			return
		}
		if r.Method == http.MethodPut {
			var v T
			if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
				ServeBadRequest(w, r)
				return
			}
			c.Put(id, v)
			return
		}
	}
	http.StripPrefix("/"+id, item).ServeHTTP(w, r)
}

func (c Collection[T]) serveRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		c.serveRootGet(w, r)
	case "POST":
		c.serveRootPost(w, r)
	default:
		ServeMethodNotAllowed(w, r)
	}
}

func (c Collection[T]) serveRootGet(w http.ResponseWriter, r *http.Request) {
	var err error
	if IsHTML(r) {
		err = c.WriteHTML(w)
	} else {
		err = json.NewEncoder(w).Encode(c)
	}
	if err != nil {
		ServeInternalServerError(w, r)
	}
}

func (c Collection[T]) serveRootPost(w http.ResponseWriter, r *http.Request) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		ServeBadRequest(w, r)
		return
	}
	id := c.Post(v)
	w.Header().Set("Location", r.URL.Path+"/"+id)
	w.WriteHeader(http.StatusCreated)
}

func (c Collection[T]) WriteHTML(w io.Writer) error {
	return collectionTmpl.Execute(w, c)
}
