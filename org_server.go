package web

// OrgServer is a server for managing organizations.
type OrgServer struct {
	// TODO
}

type Org struct {
	Members map[string]bool `json:"members"`
}
