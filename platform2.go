package web

import "strings"

type Platform2 struct {
	// Apps maps app IDs to Apps.
	Apps map[string]*App
	// Routes maps paths to app IDs.
	// Paths are r.Host + r.URL.Path.
	Routes map[string]string
}

// App returns the app for the given host and path.
func (p *Platform2) App(host, path string) (*App, bool) {
	appID, ok := p.route(host, path)
	if !ok {
		return nil, false
	}
	return p.Apps[appID], true
}

// route returns the app ID for the given host and path.
// It returns the app ID for the longest matching path.
func (p *Platform2) route(host, path string) (string, bool) {
	path = host + path
	for path != "" {
		id, ok := p.Routes[path]
		if ok {
			return id, ok
		}
		i := strings.LastIndex(path, "/")
		if i == -1 {
			break
		}
		path = path[:i]
	}
	return "", false
}
