package web

func ServeHTTPS(h http.Handler, acmeConfig ACMEConfig) error {
	server := acmeConfig.HTTPSServer(h, acmeConfig.CertManager())
	return server.ListenAndServeTLS("", "")
}
