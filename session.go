package web

type Session struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}
