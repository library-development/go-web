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

func (p *Platform) TypesDir() string {
	return filepath.Join(p.DataDir, "types")
}

func (p *Platform) FunctionsDir() string {
	return filepath.Join(p.DataDir, "functions")
}

func (p *Platform) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := BuiltRequest(r)
	if err := p.WriteRequestLog(req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Method == "POST" && r.URL.Path == "/exec" {
		var execRequest *ExecRequest
		err = json.Unmarshal(req.Body, execRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err := p.Exec(w, execRequest)
		if err != nil {
			go p.ReportError(req, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		dir := filepath.Join(p.PublicDir(), req.Host)
		http.FileServer(http.Dir(dir)).ServeHTTP(w, r)
	}
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

func (p *Platform) cmdDir() string {
	return filepath.Join(p.SourceDir(), "github.com/function-cafe/go-cmd")
}

func (p *Platform) Exec(out io.Writer, req *ExecRequest) error {
	mainFile := filepath.Join(p.cmdDir(), req.Pkg, strings.ToLower(req.Func), "main.go")
	err := golang.GenerateCLI(s.SourceDir, req.Pkg, req.Func, mainFile)
	if err != nil {
		return err
	}
	cmd := exec.Command("go", "run", mainFile)
	cmd.Dir = s.SourceDir
	cmd.Stdout = out
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	err = json.NewEncoder(stdin).Encode(req.Inputs)
	if err != nil {
		return err
	}
	err = stdin.Close()
	if err != nil {
		return err
	}
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (p *Platform) pullShopifySchemas() error {
	return nil
}

// ImportTypes reads through all the files in the /src dir and imports any exported types defined in .go files.
// The types are placed in the /types dir.
func (p *Platform) ImportTypes() error {
	return nil
}

// ImportFunctions reads through all the files in the /src dir and imports any exported functions defined in .go files.
// The functions are placed in the /functions dir.
// Methods are ignored.
func (p *Platform) ImportFunctions() error {
	return nil
}

// CloneGithubRepo clones a repo from github.com into the /src dir.
func (p *Platform) CloneGithubRepo(org, repo string) error {
	return nil
}

// CreatePublicGithubRepo creates a public repo on github.com.
func (p *Platform) CreatePublicGithubRepo(ghToken, org, repo string) error {
	return nil
}

// EditGithubRepo makes changes to a repo in /src.
// Changes are commited and published to github.com.
func (p *Platform) EditGithubRepo(org, repo string, edits []struct {
	Path  string
	Value string
}) error {
	return nil
}

// PullGithubRepo pulls a repo already in /src.
func (p *Platform) PullGithubRepo(org, repo string) error {
	return nil
}

// Generate propagates changes throughout the filesystem.
// After workers make changes to certain files, this should be run to generate dependant files.
func (p *Platform) Generate() error {
	for _, userID := range p.ListUsers() {
		p.PublishUserEdits(userID)
	}
	p.GenerateCode()
	go p.PublishCode()
	return nil
}

// ListUsers returns a list of all usernames.
func (p *Platform) ListUsers() ([]string, error) {
	return nil, nil
}
