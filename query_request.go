package web

import "lib.dev/auth"

type QueryRequest struct {
	Auth    *auth.Credentials `json:"auth"`
	Query   string            `json:"query"`
	Context string            `json:"context"`
	Input   map[string]string `json:"input"`
}
