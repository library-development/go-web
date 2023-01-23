package web

import "net/http"

type SaaS[OrgDataType http.Handler] struct {
	Orgs map[string]Org[OrgDataType]
}
