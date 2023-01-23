package web

import "net/http"

type Org[DataType any] struct {
	Data DataType
}

func (o *Org[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ServeAny(o.Data, w, r)
}
