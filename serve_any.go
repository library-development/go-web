package web

import (
	"encoding/json"
	"net/http"
	"reflect"
)

func ServeAny(v any, w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		json.NewEncoder(w).Encode(v)
		return
	}
	if r.Method == http.MethodPost {
		method := r.URL.Query().Get("method")
		t := reflect.TypeOf(v)
		for i := 0; i < t.NumMethod(); i++ {
			m := t.Method(i)
			if m.Name == method {
				m.Func.Call([]reflect.Value{reflect.ValueOf(v), reflect.ValueOf(w), reflect.ValueOf(r)})
				return
			}
		}
	}
}
