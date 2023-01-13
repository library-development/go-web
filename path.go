package web

import "strings"

type Path []string

// Parts returns the parts of the path.
func (p Path) Parts() []string {
	return []string(p)
}

// String returns the path as a string.
func (p Path) String() string {
	return "/" + strings.Join(p.Parts(), "/")
}
