package web

type LoginRequeset struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}
