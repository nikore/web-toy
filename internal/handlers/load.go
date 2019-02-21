package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"web-toy/pkg/loadgen"
)

type LoadHandler struct {
	loaders []loadgen.LoadGenerator
}

func NewLoadHandler() *LoadHandler {
	lh := &LoadHandler{
		loaders: make([]loadgen.LoadGenerator, 0, 2),
	}
	if shouldLoadCpu() {
		cpu_loader := loadgen.NewCpuLoadGenerator()
		cpu_loader.Start()
		log.Println("Started CPU loader")
		lh.loaders = append(lh.loaders, cpu_loader)
	} else {
		log.Println("No CPU load requested")
	}

	if mem_mb := desiredMemLoad(); mem_mb > 0 {
		mem_loader := loadgen.NewMemLoadGenerator(mem_mb)
		mem_loader.Start()
		log.Printf("Started memory loader (%d MiB)", mem_mb)
		lh.loaders = append(lh.loaders, mem_loader)
	} else {
		log.Println("No memory load requested")
	}
	return lh
}

func (lh *LoadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	if r.URL.Path != "/load/status" {
		// Redirect other /load/* paths until we add more functionality
		http.Redirect(w, r, "/load/status", http.StatusFound)
		return
	}

	if len(lh.loaders) > 0 {
		for _, loader := range lh.loaders {
			io.WriteString(w, loader.Status())
		}
	} else {
		io.WriteString(w, "No load generated (CPU or memory)\n")
	}
}

func shouldLoadCpu() bool {
	do := os.Getenv("WEB_TOY_LOAD_CPU")
	return (do != "" && do != "0")
}

func desiredMemLoad() int64 {
	mem_mb_str := os.Getenv("WEB_TOY_LOAD_MEM_MB")
	mem_mb, err := strconv.ParseInt(mem_mb_str, 10, 64)
	if err == nil && mem_mb > 0 {
		return mem_mb
	}
	return 0
}
