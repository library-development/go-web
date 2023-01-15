package web

// Type is the type of the file.
type Type struct {
	// IsRef is true if the value is reference.
	IsRef bool `json:"is_ref"`
	// IsList is true if the value is a list.
	IsList bool `json:"is_list"`
	// BaseType can be one of the following:
	// - any builtin Go type
	// - any type defined in the Go standard library
	// - any exported type defined in a public MIT licensed library on pkg.go.dev
	// If the ID is not builtin, it must be a valid Go import path followed by a dot and the type name.
	// For example, "github.com/library-development/go-web.File" is a valid ID.
	BaseType string `json:"base_type"`
}
