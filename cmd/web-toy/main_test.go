package main

import "testing"
import "strconv"

func TestHelloConst(t *testing.T) {
	if DefaultPort != 8080 {
		t.Error("Default port is not 8080")
	}
}

func TestGetValidatedPort(t *testing.T) {
	_, err := getValidatedPort("0")
	if err == nil {
		t.Error("getValidatedPort() didn't err on port 0")
	}
	_, err = getValidatedPort("65536")
	if err == nil {
		t.Error("getValidatedPort() didn't err on port 65536")
	}
	for port := 1; port < 65536; port++ {
		portString := strconv.Itoa(port)
		testPort, err := getValidatedPort(portString)
		if err != nil {
			t.Error("getValidatedPort() errored on valid port")
		}
		if testPort != port {
			t.Error("getValidatedPort() returned wrong port number")
		}
	}
}
