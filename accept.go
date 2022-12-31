package web

func accept(r *http.Request, accept string) bool {
	return strings.Includes(r.Header.Get("Accept"), accept)
}
