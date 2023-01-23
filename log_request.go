package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func LogRequest(r *http.Request, dir string) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	l := RequestLog{
		FromIP:  r.RemoteAddr,
		Method:  r.Method,
		Host:    r.Host,
		Path:    r.URL.Path,
		Query:   r.URL.Query(),
		Headers: r.Header,
		Body:    body,
	}
	b, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return err
	}
	timestamp := time.Now().UnixNano()
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
