package web

import (
	"net/http"
)

type Server func(p *Platform, w http.ResponseWriter, r *http.Request)
