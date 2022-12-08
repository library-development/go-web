package web

import (
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

type ACMEConfig struct {
	Hosts      []string
	AdminEmail string
	CertDir    string
}

func (c *ACMEConfig) HTTPSServer(mux http.Handler, certManager autocert.Manager) http.Server {
	return http.Server{
		Handler:   mux,
		Addr:      ":443",
		TLSConfig: certManager.TLSConfig(),
	}
}

func (c *ACMEConfig) CertManager() autocert.Manager {
	return autocert.Manager{
		Cache:      autocert.DirCache(c.CertDir),
		Prompt:     autocert.AcceptTOS,
		Email:      c.AdminEmail,
		HostPolicy: autocert.HostWhitelist(c.Hosts...),
	}
}
