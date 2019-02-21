package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func Env(w http.ResponseWriter, r *http.Request) {
	redacted := map[string]bool{
		"AWS_SECRET_ACCESS_KEY": true,
		"DATADOG_API_KEY":       true,
		"DIRENV_DIFF":           true,
		"DIRENV_WATCHES":        true,
		"GEMFURY_TOKEN":         true,
		"GITHUB_API_TOKEN":      true,
		"GOBA_SSH_ADMIN_KEY":    true,
		"NEW_RELIC_LICENSE_KEY": true,
		"ROLLBAR_CLIENT_TOKEN":  true,
		"ROLLBAR_SERVER_TOKEN":  true,
		"SECURITYSESSIONID":     true,
		"VICTOR_OPS_API_KEY":    true,
	}
	logRequest(r)

	var buffer bytes.Buffer

	buffer.WriteString("<html><head><title>Environment Variables</title></head><body>")
	for _, env := range os.Environ() {
		split := strings.SplitN(env, "=", 2)
		key, val := split[0], split[1]
		if _, banned := redacted[key]; banned {
			val = "<font color=\"red\">private</font>"
		}
		buffer.WriteString(fmt.Sprintf("<tt>%s=%s</tt><br />\n", key, val))
	}
	buffer.WriteString("</body></html>")

	serveBuffer(w, buffer)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	var buffer bytes.Buffer

	if r.URL.Path == "/" || r.URL.Path == "/hello" {
		buffer.WriteString("<html><head><title>Hello!</title></head><body><p>Hello, world!</p></body></HTML>")
		serveBuffer(w, buffer)
	} else {
		errorPath(w)
	}
}

func Ping(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	var buffer bytes.Buffer
	buffer.WriteString("pong")
	serveBuffer(w, buffer)
}

func Version(w http.ResponseWriter, r *http.Request, version string) {
	logRequest(r)
	var buffer bytes.Buffer
	buffer.WriteString(version)
	serveBuffer(w, buffer)
}

func FiveOhThree(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	http.Error(w, "503 OK", 503)
}

func Timeout(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	time.Sleep(time.Hour * 9999)
	// do nothing else.  this probably leaks, but who cares?
}

func Disconnect(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	if w == nil {
		return
	}
	wHijacker, ok := w.(http.Hijacker)
	if !ok {
		log.Println("HTTP request was not a Hijacker() interface")
		http.Error(w, "503 type conversion failed", 503)
	}
	conn, _, err := wHijacker.Hijack()
	if err != nil {
		log.Println("Couldn't hijack TCP/IP session")
		http.Error(w, "503 hijack failed", 503)
		return
	}
	conn.Close()
}

func Debug(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	data := map[string]interface{}{}

	data["method"] = r.Method
	data["headers"] = r.Header
	data["trailer"] = r.Trailer
	data["URL"] = r.URL
	data["host"] = r.Host
	data["RemoteAddr"] = r.RemoteAddr
	data["RemoteURI"] = r.RequestURI

	jsonString, err := json.Marshal(data)

	if err != nil {
		http.Error(w, "503 Error with json", 503)
		return
	}

	var buffer bytes.Buffer

	err = json.Indent(&buffer, jsonString, "", " ")

	if err != nil {
		http.Error(w, "503 Error with pretty print", 503)
		return
	}

	serveBuffer(w, buffer)
}
