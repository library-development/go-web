package web

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type AuthClient struct {
	AuthServerAddr string
}

func (c *AuthClient) CreateInviteCode() {
	panic("not implemented")
}

func (c *AuthClient) Register() {
	panic("not implemented")
}

func (c *AuthClient) Login() {
	panic("not implemented")
}

func (c *AuthClient) ValidateSession(s *Session) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	resp, err := http.Post(c.validateSessionEndpoint(), "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ErrInvalidSession
	}
	return nil
}

func (c *AuthClient) Logout() {
	panic("not implemented")
}

func (c *AuthClient) validateSessionEndpoint() string {
	return c.AuthServerAddr + "/validate-session"
}
