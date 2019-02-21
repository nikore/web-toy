package handlers

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

func errorPath(w http.ResponseWriter) {
	http.Error(w, "Path not found", 404)
}

func serveBuffer(w http.ResponseWriter, buffer bytes.Buffer) {
	io.Copy(w, &buffer)
}

func logRequest(r *http.Request) {
	log.Printf("%s: [%s] %s", r.RemoteAddr, r.Proto, r.URL.Path)
}
