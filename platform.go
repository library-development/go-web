package web

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/acme/autocert"
	"lib.dev/golang"

	"github.com/oklog/ulid/v2"
)

// Platform is the main type for the web platform.
// It implements the http.Handler interface.
// DataDir is the directory where the platform stores its data.
type Platform struct {
	// LetsEncryptEmail is the email address used to register with Let's Encrypt.
	LetsEncryptEmail string `json:"lets_encrypt_email"`

	// DataDir is a path to the directory where the platform stores its data.
	// The following subdirectories are used:
	// - src: source code
	// - logs: request logs
	// - certs: TLS certificates
	// - public: public files
	// - types: type metadata
	// - functions: function metadata
	// - bin: compiled binaries
	// - user_data: user data
	DataDir string `json:"data_dir"`

	connections *Counter
	locks       map[string]*sync.Mutex
}

func (p *Platform) db() *DB {
	locks := map[string]*sync.Mutex{}
	for lockID, lock := range p.locks {
		locks[strings.TrimPrefix(lockID, "/tables/")] = lock
	}
	return &DB{
		LocalPath: p.tablesDir(),
		Locks:     locks,
	}
}

// newID generates a new unique ID for the given scope.
func (p *Platform) newID(scope string) string {
	p.locks[scope].Lock()
	defer p.locks[scope].Unlock()
	return NewID()
}

// copyFile copies a file from one location to another.
func (p *Platform) copyFile(from, to string) error {
	data, err := os.ReadFile(from)
	if err != nil {
		return err
	}
	return os.WriteFile(to, data, os.ModePerm)
}

// Install installs the platform.
// It copies the systemd files and starts the platform service.
func (p *Platform) Install() error {
	// copy systemd files
	err := p.copyFile(filepath.Join(p.systemdDir(), "system/platform.service"), "/etc/systemd/system/platform.service")
	if err != nil {
		return err
	}

	// systemctl daemon-reload
	cmd := exec.Command("systemctl", "daemon-reload")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(output))
	}

	// systemctl start platform
	cmd = exec.Command("systemctl", "start", "platform")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(output))
	}

	return nil
}

// Upgrade upgrades the platform.
func (p *Platform) Upgrade() error {
	// Build platform
	err := p.buildCmd("platform")
	if err != nil {
		return err
	}

	// systemctl restart platform
	cmd := exec.Command("systemctl", "restart", "platform")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(output))
	}

	return nil
}

// buildCmd builds a command in /cmd.
func (p *Platform) buildCmd(name string) error {
	mainFile := filepath.Join(p.cmdDir(), name, "main.go")
	outFile := filepath.Join(p.binDir(), name)
	cmd := exec.Command("go", "build", "-o", outFile, mainFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(output))
	}
	return nil
}

// BuildAll builds all commands in /cmd.
func (p *Platform) BuildAll() error {
	entries, err := os.ReadDir(p.cmdDir())
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			if err := p.buildCmd(entry.Name()); err != nil {
				return err
			}
		}
	}
	return nil
}

// autocertManager builds the autocert.Manager for internal use.
func (p *Platform) autocertManager() autocert.Manager {
	return autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache(p.certDir()),
		HostPolicy: func(_ context.Context, host string) error {
			hosts, err := p.listHosts()
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

// listHosts returns a list of strings containing the names of all hosts on the platform.
func (p *Platform) listHosts() ([]string, error) {
	hosts := []string{}
	entries, err := os.ReadDir(p.publicDir())
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

// Start starts the web platform.
func (p *Platform) Start() error {
	manager := p.autocertManager()
	return http.Serve(manager.Listener(), p)
}

// /systemd
func (p *Platform) systemdDir() string {
	return filepath.Join(p.DataDir, "systemd")
}

// /src
func (p *Platform) sourceDir() string {
	return filepath.Join(p.DataDir, "src")
}

// /logs
func (p *Platform) logsDir() string {
	return filepath.Join(p.DataDir, "logs")
}

// /certs
func (p *Platform) certDir() string {
	return filepath.Join(p.DataDir, "certs")
}

// /public
func (p *Platform) publicDir() string {
	return filepath.Join(p.DataDir, "public")
}

// /types
func (p *Platform) typesDir() string {
	return filepath.Join(p.DataDir, "types")
}

// /functions
func (p *Platform) functionsDir() string {
	return filepath.Join(p.DataDir, "functions")
}

// /bin
func (p *Platform) binDir() string {
	return filepath.Join(p.DataDir, "bin")
}

// /files
func (p *Platform) filesDir() string {
	return filepath.Join(p.DataDir, "files")
}

// fileDir is a nicer name than filesDir.
func (p *Platform) fileDir() string {
	return p.filesDir()
}

// /auth
func (p *Platform) authDir() string {
	return filepath.Join(p.DataDir, "auth")
}

// /errors
func (p *Platform) errorsDir() string {
	return filepath.Join(p.DataDir, "errors")
}

// logRequest logs the request to the logs directory.
func (p *Platform) logRequest(r *http.Request) error {
	req, err := BuildRequest(r)
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	var logFile string
	filename := req.Timestamp
	for {
		logFile = filepath.Join(p.logsDir(), strconv.Itoa(int(filename)))
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			break
		}
		filename++
	}
	return os.WriteFile(logFile, b, os.ModePerm)
}

// reverseProxyAddress returns the address of the reverse proxy for the request and true if the request should be proxied.
func (p *Platform) reverseProxyAddress(r *http.Request) (string, bool) {
	dirPath := filepath.Join(p.publicDir(), r.Host, r.URL.Path)
	for {
		filePath := filepath.Join(dirPath, "REVERSE_PROXY")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			dirPath = filepath.Dir(dirPath)
			if dirPath == p.publicDir() {
				return "", false
			}
		} else {
			b, err := os.ReadFile(filePath)
			if err != nil {
				panic(err)
			}
			return string(b), true
		}
	}
}

// htmlTemplate returns an html.Template for browser requests.
func (p *Platform) htmlTemplate() *template.Template {
	return template.Must(template.New("html").Parse(htmlTemplate))
}

// getFile returns the file with the given path.
// If the file does not exist, an error is returned.
// If the file is a directory, /index is appended to the path.
func (p *Platform) getFile(path string) (*File, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		path = filepath.Join(path, "index")
	}
	var file *File
	b, err := os.ReadFile(filepath.Join(p.fileDir(), path))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, file)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// handleError writes an error to the response if err is not nil.
func (p *Platform) handleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	if os.IsNotExist(err) {
		if isHTML(r) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte("404 - Not Found"))
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("\"meta\":{\"error\":\"not found\"}}"))
		}
		return
	}

	p.reportError(r, err)

	if isHTML(r) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("500 - Internal Server Error"))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("\"meta\":{\"error\":\"internal\"}}"))
	}
}

func (p *Platform) tablesDir() string {
	return filepath.Join(p.DataDir, "tables")
}

func (p *Platform) tableType(tableID string) (*golang.Ident, error) {
	path := filepath.Join(p.tablesDir(), tableID, "type")
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var t *golang.Ident
	err = json.Unmarshal(b, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (p *Platform) table(id string) (*Table, error) {
	return &Table{
		LocalPath: filepath.Join(p.tablesDir(), id),
	}, nil
}

func (p *Platform) createTable(id string, t *golang.Ident) error {
	path := filepath.Join(p.tablesDir(), id, "type")
	b, err := json.Marshal(t)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, b, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// ServeHTTPNext is the next version of the ServeHTTPNext function.
// It is a work in progress.
func (p *Platform) ServeHTTPNext(w http.ResponseWriter, r *http.Request) {
	p.connections.Inc()

	err := p.logRequest(r)
	if err != nil {
		p.reportError(r, err)
	}

	revProxyAddr, isRevProxy := p.reverseProxyAddress(r)
	if isRevProxy {
		proxyURL, err := url.Parse(revProxyAddr)
		if err != nil {
			p.reportError(r, err)
			return
		}

		proxy := &httputil.ReverseProxy{
			Director: func(r *http.Request) {
				r.URL.Scheme = proxyURL.Scheme
				r.URL.Host = proxyURL.Host
				r.URL.Path = filepath.Join(proxyURL.Path, r.URL.Path)
			},
		}
		proxy.ServeHTTP(w, r)
		return
	}

	// main logic
	path := filepath.Join(r.Host, r.URL.Path)
	file, err := p.getFile(path)
	if err != nil {
		p.handleError(w, r, err)
		return
	}
	if isGET(r) {
		if isHTML(r) {
			p.htmlTemplate().Execute(w, file)
		} else {
			json.NewEncoder(w).Encode(file)
		}
	}
	if isPOST(r) {
	}

	p.connections.Dec()
}

// ServeHTTP is the main request handler for the platform.
func (p *Platform) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := p.logRequest(r)
	if err != nil {
		p.reportError(r, err)
	}

	revProxyAddr, isRevProxy := p.reverseProxyAddress(r)
	if isRevProxy {
		proxyURL, err := url.Parse(revProxyAddr)
		if err != nil {
			p.reportError(r, err)
			return
		}

		proxy := &httputil.ReverseProxy{
			Director: func(r *http.Request) {
				r.URL.Scheme = proxyURL.Scheme
				r.URL.Host = proxyURL.Host
				r.URL.Path = filepath.Join(proxyURL.Path, r.URL.Path)
			},
		}
		proxy.ServeHTTP(w, r)
		return
	}

	dir := filepath.Join(p.publicDir(), r.Host)
	http.FileServer(http.Dir(dir)).ServeHTTP(w, r)
}

// isAuthorized returns true if the request is authorized to continue.
func (p *Platform) isAuthorized(r *http.Request, file *File) bool {
	if r.Method == http.MethodGet {
		return true
	}
	sessionToken := r.Header.Get("Authorization")
	userID := p.UserID(sessionToken)
	_, ok := file.Metadata.Owners[userID]
	return ok
}

// UserID returns the user ID for the given session token.
func (p *Platform) UserID(sessionToken string) string {
	sessionFile := filepath.Join(p.authDir(), "sessions", sessionToken)
	userID, _ := os.ReadFile(sessionFile)
	return string(userID)
}

// WriteFile writes the file for the given host and path.
func (p *Platform) WriteFile(host, path string, file *File) error {
	fp := filepath.Join(p.filesDir(), host, path)
	err := os.MkdirAll(filepath.Dir(fp), os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	err = enc.Encode(file)
	if err != nil {
		return err
	}
	return nil
}

// readFile returns the file for the given host and path.
func (p *Platform) readFile(host, path string) (*File, error) {
	f, err := os.Open(filepath.Join(p.filesDir(), host, path))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	var file File
	err = dec.Decode(&file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// newName generates a new random string for a file name.
func (p *Platform) newName() string {
	return ulid.Make().String()
}

// reportError writes the error to file.
func (p *Platform) reportError(r *http.Request, err error) {
	if err == nil {
		panic(err)
	}

	b, err := json.Marshal(err)
	if err != nil {
		panic(err)
	}

	createdAt := time.Now().UTC().UnixNano()

	f := &File{
		Metadata: Metadata{
			Type: "error",
			Owners: map[string]bool{
				"mikerybka": true,
			},
			Public:    false,
			Name:      p.newName(),
			CreatedAt: createdAt,
		},
		Data: b,
	}

	p.WriteFile("mikerybka.com", "/feedback/platform/errors/"+strconv.Itoa(int(createdAt)), f)
}

func (p *Platform) funcCmdPkg() string {
	return "github.com/function-cafe/go-cmd"
}

func (p *Platform) methodCmdPkg() string {
	return "github.com/webmachine-dev/go-cmd"
}

func (p *Platform) funcCmdDir() string {
	return filepath.Join(p.sourceDir(), p.funcCmdPkg())
}

func (p *Platform) cmdDir() string {
	return filepath.Join(p.DataDir, "cmd")
}

func (p *Platform) methodCmdDir() string {
	return filepath.Join(p.sourceDir(), p.methodCmdPkg())
}

func (p *Platform) goFuncMainFile(pkg, fn string) string {
	return filepath.Join(p.funcCmdDir(), pkg, strings.ToLower(fn), "main.go")
}

func (p *Platform) goMethodMainFile(pkg, fn string) string {
	return filepath.Join(p.methodCmdDir(), pkg, strings.ToLower(fn), "main.go")
}

func (p *Platform) goFuncBinFile(pkg, fn string) string {
	return filepath.Join(p.binDir(), p.funcCmdPkg(), pkg, strings.ToLower(fn))
}

func (p *Platform) goMethodBinFile(pkg, fn string) string {
	return filepath.Join(p.binDir(), p.methodCmdPkg(), pkg, strings.ToLower(fn))
}

func (p *Platform) generateGoFuncMainFile(pkg, fn string) error {
	mainFile := p.goFuncMainFile(pkg, fn)
	err := golang.GenerateCLI(p.sourceDir(), pkg, fn, mainFile)
	if err != nil {
		return err
	}
	return nil
}

func (p *Platform) generateGoMethodMainFile(pkg, fn string) error {
	mainFile := p.goMethodMainFile(pkg, fn)
	err := golang.GenerateCLI(p.sourceDir(), pkg, fn, mainFile)
	if err != nil {
		return err
	}
	return nil
}

func (p *Platform) BuildGoFunc(pkg, fn string) error {
	err := p.generateGoFuncMainFile(pkg, fn)
	if err != nil {
		return err
	}
	mainFile := p.goFuncMainFile(pkg, fn)
	binFile := p.goFuncBinFile(pkg, fn)
	cmd := exec.Command("go", "build", "-o", binFile, mainFile)
	cmd.Dir = p.DataDir
	return cmd.Run()
}

func (p *Platform) exec(out io.Writer, req *ExecRequest) error {
	cmd := exec.Command(p.goFuncMainFile(req.Pkg, req.Func))
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

// importTypes reads through all the files in the /src dir and imports any exported types defined in .go files.
// The types are placed in the /types dir.
func (p *Platform) importTypes() error {
	return nil
}

// importFunctions reads through all the files in the /src dir and imports any exported functions defined in .go files.
// The functions are placed in the /functions dir.
// Methods are ignored.
func (p *Platform) importFunctions() error {
	return nil
}

// cloneGithubRepo clones a repo from github.com into the /src dir.
func (p *Platform) cloneGithubRepo(org, repo string) error {
	return nil
}

// createPublicGithubRepo creates a public repo on github.com.
func (p *Platform) createPublicGithubRepo(ghToken, org, repo string) error {
	return nil
}

// editGithubRepo makes changes to a repo in /src.
// Changes are commited and published to github.com.
func (p *Platform) editGithubRepo(org, repo string, edits []struct {
	Path  string
	Value string
}) error {
	return nil
}

// pullGithubRepo pulls a repo already in /src.
func (p *Platform) pullGithubRepo(org, repo string) error {
	return nil
}

// listUsers returns a list of all usernames.
func (p *Platform) listUsers() ([]string, error) {
	return nil, nil
}

// publishUserEdits copies public files owned by owner from thier data dir to /public.
func (p *Platform) publishUserEdits(owner string) error {
	return nil
}

// generateCode generates code from the /src dir.
func (p *Platform) generateCode() error {
	return nil
}

// publishCode pushes code to github.com.
func (p *Platform) publishCode() error {
	return nil
}

// generateTypes generates types from the /src dir.
func (p *Platform) generateTypes() error {
	return nil
}

// generateSchemas generates schemas from the /src dir.
func (p *Platform) generateSchemas() error {
	return nil
}

// generateFunctions generates functions from the /src dir.
func (p *Platform) generateFunctions() error {
	return nil
}

// generatePublic generates the rest of the public dir.
func (p *Platform) generatePublic() error {
	return nil
}

// htmlViewTemplate returns the HTML template for viewing the given type.
func (p *Platform) htmlViewTemplate(t string) *template.Template {
	templatePath := filepath.Join(p.typesDir(), t, "view.html")
	tmpl, _ := template.New("view.html").ParseFiles(templatePath)
	return tmpl
}

func (p *Platform) getType(id golang.Ident) (*golang.Type, error) {
	return nil, nil
}
