package web

import "github.com/library-development/go-english"

// Metadata is the metadata for a file.
type Metadata struct {
	// Type is the type of the file.
	// It can be one of the following:
	// - any builtin Go type
	// - any type defined in the Go standard library
	// - any exported type defined in a public MIT licensed library on pkg.go.dev
	Type string `json:"type"`
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
