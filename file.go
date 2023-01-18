package web

import (
	"encoding/json"
	"net/http"

	"github.com/library-development/go-english"
)

// File is core concept in the web library.
// It is used throughout the package to represent files and directories.
// It doesn't actually contain the file's contents, but a SHA256 hash of the file's contents.
type File struct {
	// Type is the type of the file.
	Type Type `json:"type"`
	// Owner is the ID of the organization that owns the file.
	Owner string `json:"owner"`
	// Owners is a map of user IDs that own the file.
	// Deprecated: Use Owner instead.
	Owners map[string]bool `json:"owners"`
	// Name is the name of the file.
	// Deprecated: Use EnglishName instead.
	Name string `json:"name"`
	// EnglishName is the English-language name of the file.
	EnglishName english.Name `json:"english_name"`
	// Doc is the documentation for the file.
	Doc string `json:"doc"`
	// Comments are a list of comments on the file.
	Comments []Comment `json:"comments"`
	// Public decides whether the file is publicly available or not.
	// This means it can be viewed by anyone.
	// It can still only be edited by members of the owner organization.
	Public bool `json:"public"`
	// CreatedAt is the Unix timestamp of when the file was created.
	CreatedAt int64 `json:"created_at"`
	// UpdatedAt is the Unix timestamp of when the file was last updated.
	UpdatedAt int64 `json:"updated_at"`
	// SHA256 is the SHA256 hash of the file.
	SHA256 string `json:"sha256"`
	// Size is the size of the file in bytes.
	Size int64 `json:"size"`
}

func (f *File) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
