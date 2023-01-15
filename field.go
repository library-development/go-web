package web

import "lib.dev/english"

type Field struct {
	Type        Type         `json:"type"`
	EnglishName english.Name `json:"name"`
}
