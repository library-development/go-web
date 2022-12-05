package web

type Model struct {
	Name    string     `json:"name"`
	Fields  []Field    `json:"fields"`
	Methods []Function `json:"methods"`
}
