package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/soheilhy/cmux"
	"web-toy/internal/grpc"
	"web-toy/internal/handlers"
	"web-toy/pkg/junkport"
)

var version string = "UNDEFINED"

const DefaultPort = 8080

func getValidatedPort(portString string) (int, error) {
	port, err := strconv.Atoi(portString)
	if err == nil && port >= 1 && port <= 65535 {
		return port, nil
	}
	if err != nil {
		return 0, err
	}
	return 0, errors.New("Port out of range (1 - 65535)")
}

func httpServe(lis net.Listener) error {
	fiveOhThreeDelay, err := strconv.Atoi(os.Getenv("FIVEOHTHREE_FUSE_SECS"))
	if err != nil {
		fiveOhThreeDelay = 0
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.Hello)
	mux.HandleFunc("/ping", handlers.Ping)
	mux.HandleFunc("/disconnect", handlers.Disconnect)
	mux.HandleFunc("/env", handlers.Env)
	mux.HandleFunc("/timeout", handlers.Timeout)
	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) { handlers.Version(w, r, version) })
	mux.HandleFunc("/503/", handlers.FiveOhThree)
	mux.HandleFunc("/debug", handlers.Debug)

	healthHandler := handlers.NewHealthHandler(fiveOhThreeDelay)
	mux.Handle("/health", healthHandler)

	loadHandler := handlers.NewLoadHandler()
	mux.Handle("/load/", loadHandler)

	s := &http.Server{Handler: mux}
	return s.Serve(lis)
}

func main() {
	port := 0
	if portEnv, exists := os.LookupEnv("HTTP_PORT"); exists {
		httpPort, err := getValidatedPort(portEnv)
		if err != nil {
			log.Fatalf("Invalid port in $HTTP_PORT: %s", err.Error())
		}
		port = httpPort
	}
	if port == 0 {
		log.Println("$HTTP_PORT not defined.  Listening on default port")
		port = DefaultPort
	}

	badPort, err := getValidatedPort(os.Getenv("JUNK_PORT"))
	if err == nil {
		go func() {
			bpErr := junkport.Serve(badPort)
			// Should never return
			log.Fatalf("junk-port server exited with error: %s", bpErr.Error())
		}()
	} else {
		log.Printf("$JUNK_PORT not defined. Listening on default port")
		go func() {
			bpErr := junkport.Serve(1332)
			// Should never return
			log.Fatalf("junk-port server exited with error: %s", bpErr.Error())
		}()
	}

	log.Printf("Starting web-toy on port %d", port)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal(err)
	}

	m := cmux.New(listener)

	grpcListener := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpListener := m.Match(cmux.HTTP1Fast())

	go httpServe(httpListener)
	go grpc.Serve(grpcListener)

	m.Serve()
}
