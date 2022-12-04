package web

import "github.com/library-development/go-auth"

type QueryRequest struct {
	Auth    *auth.Credentials `json:"auth"`
	Query   string            `json:"query"`
	Context string            `json:"context"`
	Input   map[string]string `json:"input"`
}
