package web

type CreateResponse struct {
	ID    string `json:"id"`
	Error string `json:"error",omitempty,omitemptykey:""`
}
