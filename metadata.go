package web

import "lib.dev/golang"

type Metadata struct {
	Type     golang.Ident    `json:"type"`
	Owners   map[string]bool `json:"owners"`
	Name     string          `json:"name"`
	Doc      string          `json:"doc"`
	Comments []Comment       `json:"comments"`
	Public   bool            `json:"public"`
}
