package junkport

import "errors"
import "log"
import "os/exec"
import "net"

const maxErrors int = 5
const fortuneBin string = "/usr/games/fortune"

// Serve a TCP/IP port with something besides valid HTTP.
func Serve(port int) error {
	addr := net.TCPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: port,
	}
	listener, err := net.ListenTCP("tcp4", &addr)
	if err != nil {
		return err
	}
	log.Printf("Junk-port server listening on %s", addr.String())

	errorCount := 0
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			errorCount += 1
			if errorCount >= maxErrors {
				return errors.New("ERROR: junk-port server: accept() error count threshold exceeded")
			}
			log.Printf("WARNING: junk-port server: accept() failed: %s", err.Error())
			continue
		}
		errorCount = 0
		log.Printf("junk-port server: Got connection from %s", conn.RemoteAddr().String())
		if err := sendJunkMessage(conn); err != nil {
			log.Printf("WARNING: junk-port server: sendJunkMessage() failed: %s", err.Error())
		}
		if err := conn.Close(); err != nil {
			log.Printf("WARNING: junk-port server: sendJunkMessage(): close() failed: %s", err.Error())
		}
	}

	return errors.New("ERROR: junk-port server exited unexpectedly")
}

func sendJunkMessage(conn *net.TCPConn) error {
	fortune, err := exec.Command(fortuneBin).Output()
	if err != nil {
		return err
	}
	outBytes, err := conn.Write(fortune)
	if err != nil {
		return err
	}
	if outBytes != len(fortune) {
		log.Printf("WARNING: junk-port server: sendJunkMessage(): wrote %d bytes of %d-byte string", outBytes, len(fortune))
	}
	return nil
}
