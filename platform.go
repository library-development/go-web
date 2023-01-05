package web

import (
	"bytes"
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

	"golang.org/x/crypto/acme/autocert"
	"lib.dev/golang"
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
}

func (p *Platform) copyFile(from, to string) error {
	data, err := os.ReadFile(from)
	if err != nil {
		return err
	}
	return os.WriteFile(to, data, os.ModePerm)
}

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

func (p *Platform) Update() error {
	// Build platform
	err := p.BuildCmd("platform")
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

func (p *Platform) BuildCmd(name string) error {
	mainFile := filepath.Join(p.cmdDir(), name, "main.go")
	outFile := filepath.Join(p.binDir(), name)
	cmd := exec.Command("go", "build", "-o", outFile, mainFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(output))
	}
	return nil
}

func (p *Platform) BuildAll() error {
	entries, err := os.ReadDir(p.cmdDir())
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			if err := p.BuildCmd(entry.Name()); err != nil {
				return err
			}
		}
	}
	return nil
}

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

func (p *Platform) systemdDir() string {
	return filepath.Join(p.DataDir, "systemd")
}

func (p *Platform) sourceDir() string {
	return filepath.Join(p.DataDir, "src")
}

func (p *Platform) logsDir() string {
	return filepath.Join(p.DataDir, "logs")
}

func (p *Platform) certDir() string {
	return filepath.Join(p.DataDir, "certs")
}

func (p *Platform) publicDir() string {
	return filepath.Join(p.DataDir, "public")
}

func (p *Platform) typesDir() string {
	return filepath.Join(p.DataDir, "types")
}

func (p *Platform) functionsDir() string {
	return filepath.Join(p.DataDir, "functions")
}

func (p *Platform) binDir() string {
	return filepath.Join(p.DataDir, "bin")
}

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

// Proxy returns the port of the proxy for the given host.
func (p *Platform) proxy(host string) (string, bool) {
	proxyFile := filepath.Join(p.proxyDir(), host)
	port, err := os.ReadFile(proxyFile)
	if err != nil {
		return "", false
	}
	return string(port), true
}

// proxyDir returns the directory where the proxy files are stored.
func (p *Platform) proxyDir() string {
	return filepath.Join(p.DataDir, "proxy")
}

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

func (p *Platform) htmlTemplate() *template.Template {
	return template.Must(template.New("html").Parse(htmlTemplate))
}

func (p *Platform) fileDir() string {
	return filepath.Join(p.DataDir, "files")
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

func (p *Platform) writeError(w http.ResponseWriter, r *http.Request, err error, code int) {
	w.WriteHeader(code)
	if isHTML(r) {
		fmt.Fprintf(w, "Error: %v", err)
	} else {
		json.NewEncoder(w).Encode(err)
	}
}

func (p *Platform) handleError(w http.ResponseWriter, r *http.Request, err error) {
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

	p.ReportError(r, err)

	if isHTML(r) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("500 - Internal Server Error"))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("\"meta\":{\"error\":\"internal\"}}"))
	}
}

func (p *Platform) GenericServe(w http.ResponseWriter, r *http.Request) {
	p.connections.Inc()
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

func (p *Platform) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := p.logRequest(r)
	if err != nil {
		p.ReportError(r, err)
	}

	revProxyAddr, isRevProxy := p.reverseProxyAddress(r)
	if isRevProxy {
		proxyURL, err := url.Parse(revProxyAddr)
		if err != nil {
			p.ReportError(r, err)
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

// authDir returns the directory where the auth files are stored.
func (p *Platform) authDir() string {
	return filepath.Join(p.DataDir, "auth")
}

// IsAuthorized returns true if the request is authorized to continue.
func (p *Platform) IsAuthorized(r *http.Request, file *File) bool {
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
	f, err := os.Create(filepath.Join(p.filesDir(), host, path))
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

func (p *Platform) filesDir() string {
	return filepath.Join(p.DataDir, "files")
}

// ReadFile returns the file for the given host and path.
func (p *Platform) ReadFile(host, path string) (*File, error) {
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

func (p *Platform) handleGET(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(filepath.Join(p.publicDir(), r.Host, r.URL.Path))
	if err != nil {
		if os.IsNotExist(err) {
			p.handleNotFound(w, r)
			return
		}
		p.ReportError(r, err)
		return
	}
	defer f.Close()
	t := p.Type(r.Host, r.URL.Path)
	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		p.htmlHeader(w, r)
		p.writeHTML(w, t, f)
		p.htmlFooter(w, r)
	} else {
		var metadata struct {
			Type  string `json:"type"`
			ID    string `json:"id"`
			Error error  `json:"error"`
		}
		metadata.Type = t
		metadata.ID = r.Host + r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"meta\":"))
		enc := json.NewEncoder(w)
		enc.Encode(metadata)
		w.Write([]byte(",\"data\":"))
		io.Copy(w, f)
		w.Write([]byte("}"))
	}
}

func (p *Platform) htmlHeader(w http.ResponseWriter, r *http.Request) {
}

func (p *Platform) htmlFooter(w http.ResponseWriter, r *http.Request) {
}

// writeHTML converts the json input d to html as type t and writes it to w.
func (p *Platform) writeHTML(w io.Writer, t string, d io.Reader) {
	switch t {
	default:
		fmt.Fprintf(w, "<pre>%s</pre>", d)
	}
}

func (p *Platform) handleNotFound(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("404 - Not Found"))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("\"meta\":{\"error\":\"not found\"}}"))
	}
}

func (p *Platform) handlePOST(w http.ResponseWriter, r *http.Request) {
	req := &ExecRequest{
		Token: r.Header.Get("Token"),
		Pkg:   filepath.Dir(p.funcCmdPkg() + "/" + r.Host + r.URL.Path),
		Func:  filepath.Base(r.URL.Path),
	}
	json.NewDecoder(r.Body).Decode(&req.Inputs)
	out := &bytes.Buffer{}
	err := p.exec(out, req)
	if err != nil {
		p.ReportError(r, err)
	}
	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		w.Header().Set("Content-Type", "text/html")
		p.htmlHeader(w, r)
		p.writeHTML(w, "json", out)
		p.htmlFooter(w, r)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(out.Bytes())
	}
}

func (p *Platform) ReportError(r *http.Request, err error) {
	fmt.Println("Error:", err)
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

// ListUsers returns a list of all usernames.
func (p *Platform) ListUsers() ([]string, error) {
	return nil, nil
}

// PublishUserEdits copies public files owned by owner from thier data dir to /public.
func (p *Platform) PublishUserEdits(owner string) error {
	return nil
}

// GenerateCode generates code from the /src dir.
func (p *Platform) GenerateCode() error {
	return nil
}

// PublishCode pushes code to github.com.
func (p *Platform) PublishCode() error {
	return nil
}

// GenerateTypes generates types from the /src dir.
func (p *Platform) GenerateTypes() error {
	return nil
}

// GenerateSchemas generates schemas from the /src dir.
func (p *Platform) GenerateSchemas() error {
	return nil
}

// GenerateFunctions generates functions from the /src dir.
func (p *Platform) GenerateFunctions() error {
	return nil
}

// GeneratePublic generates the rest of the public dir.
func (p *Platform) GeneratePublic() error {
	return nil
}

// Type returns the type at a given path.
func (p *Platform) Type(host, path string) string {
	fileinfo, err := os.Stat(filepath.Join(p.publicDir(), host, path))
	if err != nil {
		return ""
	}
	if fileinfo.IsDir() {
		return "folder"
	}
	parts := PathParts(path)
	switch host {
	case "api.schema.cafe":
		if len(parts) == 0 {
			return "folder"
		}
		if parts[0] == "schemas" {
			return "schema"
		}
	default:
		err := fmt.Errorf("unknown host: %s", host)
		p.ReportError(nil, err)
		return ""
	}
	return ""
}

// HTMLViewTemplate returns the HTML template for viewing the given type.
func (p *Platform) HTMLViewTemplate(t string) *template.Template {
	templatePath := filepath.Join(p.typesDir(), t, "view.html")
	tmpl, _ := template.New("view.html").ParseFiles(templatePath)
	return tmpl
}
