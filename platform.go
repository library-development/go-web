package web

import (
	"context"
	"fmt"
	"path/filepath"

	"golang.org/x/crypto/acme/autocert"
)

type Platform struct {
	DataDir          string
	LetsEncryptEmail string
}

func (p *Platform) autocertManager() autocert.Manager {
	return autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache(p.CertDir()),
		HostPolicy: func(_ context.Context, host string) error {
			hosts, err := p.listHosts(dir)
			if err != nil {
				return err
			}
			for _, h := range hosts {
				if h == host {
					return nil
				}
			}
			return fmt.Errorf("host %q not allowed", host)
		},
		Email: p.LetsEncryptEmail,
	}
}

func (p *Platform) listHosts() ([]string, error) {
	hosts := []string{}
	entries, err := os.ReadDir(p.PublicDir())
	if err != nil {
		return hosts, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			hosts = append(hosts, entry.Name())
		}
	}
	return hosts, nil
}

func (p *Platform) Start() error {
	return http.Serve(p.autocertManager().Listener(), p)
}

func (p *Platform) SourceDir() string {
	return filepath.Join(p.DataDir, "src")
}

func (p *Platform) LogsDir() string {
	return filepath.Join(p.DataDir, "logs")
}

func (p *Platform) CertDir() string {
	return filepath.Join(p.DataDir, "certs")
}

func (p *Platform) PublicDir() string {
	return filepath.Join(p.DataDir, "public")
}

func (p *Platform) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := BuiltRequest(r)
	if err := p.WriteRequestLog(req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Method == "POST" && r.URL.Path == "/exec" {
		var execRequest ExecRequest
		err = json.Unmarshal(req.Body, &execRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		result, err := p.Exec(&execRequest)
		if err != nil {
			go p.ReportError(req, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			go p.ReportError(req, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	dir := filepath.Join(p.PublicDir(), req.Host)
	http.FileServer(http.Dir(dir)).ServeHTTP(w, r)
}

func (p *Platform) ReportError(req *Request, err error) {
	fmt.Println("Error:", err)
}

func (p *Platform) WriteRequestLog(r *Request) error {
	timestamp := time.Now().UnixNano()
	b, err := json.MarshalIndent(BuiltRequest(r), "", "  ")
	if err != nil {
		return err
	}
	var logFile string
	for {
		logFile = filepath.Join(dir, strconv.Itoa(int(timestamp)))
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			break
		}
		timestamp++
	}
	return os.WriteFile(logFile, b, os.ModePerm)
}

func (p *Platform) Exec(req *ExecRequest) (any, error) {
	return nil, nil
}
