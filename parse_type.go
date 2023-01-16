package web

import "strings"

func ParseType(s string) Type {
	t := Type{}
	if strings.HasPrefix(s, "[]") {
		t.IsList = true
		s = strings.TrimPrefix(s, "[]")
	}
	if strings.HasPrefix(s, "*") {
		t.IsRef = true
		s = strings.TrimPrefix(s, "*")
	}
	t.BaseType = s
	return t
}
