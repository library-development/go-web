package web

type Metadata struct {
	Type     golang.Ident    `json:"type"`
	Owners   map[string]bool `json:"owners"`
	Name     string          `json:"name"`
	Doc      string          `json:"doc"`
	Comments []Comment       `json:"comments"`
}

// WriteHTML writes the HTML representation of the metadata to the writer.
func (m *Metadata) WriteHTML(w io.Writer) error {
	return nil
}
