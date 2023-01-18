package web

import "github.com/library-development/go-english"

type Field struct {
	Type        Type         `json:"type"`
	EnglishName english.Name `json:"name"`
}
