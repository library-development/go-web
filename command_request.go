package web

import "github.com/library-development/go-auth"

type CommandRequest struct {
	Auth    *auth.Credentials `json:"auth"`
	Command string            `json:"command"`
	Context string            `json:"context"`
	Input   map[string]string `json:"input"`
}
