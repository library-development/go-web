package web

import (
	"encoding/json"

	"github.com/library-development/go-english"
)

type HomePage struct {
	Title string
}

func (h HomePage) Type() Type {
	return Type{
		BaseType: "github.com/library-development/go-web/HomePage",
	}
}

func (h HomePage) File() File {
	b, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}

	return File{
		Type:        h.Type(),
		Owner:       "internal",
		Public:      true,
		EnglishName: english.ParseName(h.Title),
		SHA256:      checksum(b),
		Size:        int64(len(b)),
	}
}
