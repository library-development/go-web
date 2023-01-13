package web

import "lib.dev/english"

// Org is an organization.
// Organizations are
type Organization struct {
	// EnglishName is the English-language name of the organization.
	EnglishName english.Name `json:"english_name"`
	// Memebers is a list of UserIDs of the members of the organization.
	Members []string `json:"members"`
}
