package web

import (
	"io"

	"github.com/library-development/go-english"
)

type Resource interface {
	Type() string
	WriteJSON(w io.Writer) error
	WriteHTML(w io.Writer) error
	CallMethod(method string, input io.Reader, output io.Writer) error
	ID() string
	Name() english.Name
	Doc() string
}
