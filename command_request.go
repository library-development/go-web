package web

import "lib.dev/auth"

type CommandRequest struct {
	Auth    *auth.Credentials `json:"auth"`
	Command string            `json:"command"`
	Context string            `json:"context"`
	Input   map[string]string `json:"input"`
}
