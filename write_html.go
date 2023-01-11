package web

import (
	"encoding/json"
	"fmt"
	"io"

	"lib.dev/golang"
)

func writeHTML(t golang.Ident, w io.Writer, b []byte) error {
	switch t.From {
	case "":
		switch t.Name {
		case "string":
			var d string
			if err := json.Unmarshal(b, &d); err != nil {
				return err
			}
			return writeHTMLString(w, d)
		case "int":
			var d int
			if err := json.Unmarshal(b, &d); err != nil {
				return err
			}
			return writeHTMLString(w, fmt.Sprintf("%d", d))
		case "bool":
			var d bool
			if err := json.Unmarshal(b, &d); err != nil {
				return err
			}
			return writeHTMLString(w, fmt.Sprintf("%t", d))
		}
	}
	return fmt.Errorf("cannot write %s as HTML", t)
}

func writeHTMLString(w io.Writer, s string) error {
	_, err := w.Write([]byte(s))
	return err
}
