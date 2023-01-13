package web

import "strings"

// ParsePath parses a url.URL.Path into a web.Path.
func ParsePath(p string) Path {
	var path Path
	if p == "" {
		return path
	}
	if p == "/" {
		return path
	}
	for _, part := range strings.Split(p[1:], "/") {
		path = append(path, part)
	}
	return path
}
