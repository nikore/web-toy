package handlers

import (
	"io"
	"log"
	"net/http"
	"time"
)

type HealthHandler struct {
	fiveOhThreeFuseEnabled bool
	fiveOhThreeFuseTime    time.Time
}

func NewHealthHandler(fuseSeconds int) *HealthHandler {
	h := new(HealthHandler)
	if fuseSeconds > 0 {
		log.Printf("WARNING: 503 fuse timer on /health enabled: %d second(s) from start", fuseSeconds)
		h.fiveOhThreeFuseEnabled = true
		h.fiveOhThreeFuseTime = time.Now().Add(time.Duration(fuseSeconds) * time.Second)
	}

	return h
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	if h.fiveOhThreeFuseEnabled && time.Now().After(h.fiveOhThreeFuseTime) {
		log.Print("Fuse time passed!  Serving 503 on /health")
		http.Error(w, "503 fuse time exceeded", 503)
	} else {
		io.WriteString(w, "OK doc")
	}
}
