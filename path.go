package web

import (
	"strings"
)

type Path []string

// Parts returns the parts of the path.
func (p Path) Parts() []string {
	return []string(p)
}

// String returns the path as a string.
func (p Path) String() string {
	return "/" + strings.Join(p.Parts(), "/")
}

func (p Path) Length() int {
	return len(p)
}

func (p Path) Last() string {
	if p.Length() < 1 {
		return ""
	}
	return p[len(p)-1]
}

func (p Path) Append(name string) Path {
	return append(p, name)
}

func (p Path) SecondLast() string {
	if p.Length() < 2 {
		return ""
	}
	return p[len(p)-2]
}

func (p Path) Pop() Path {
	return p[:len(p)-1]
}
