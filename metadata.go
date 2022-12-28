package web

type Metadata struct {
	Type     golang.Ident `json:"type"`
	Name     string       `json:"name"`
	Doc      string       `json:"doc"`
	Comments []Comment    `json:"comments"`
}
