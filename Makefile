VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
			cat $(CURDIR)/.version 2> /dev/null || echo v0.1)
PACKAGES = $$(go list ./... | grep -v '/vendor/')
FILES = $(shell find . -type f -name '*.go' -print)
TGT=bin/web-toy
CLIENT=bin/web-toy-client

export GO111MODULE=on

.PHONY: default all build fmt test clean

default: fmt test build

build: $(TGT) $(CLIENT)

$(TGT):
	go build -v \
	  -tags release \
		-ldflags '-X main.Version=$(VERSION)' \
		-o $(TGT) \
  	cmd/web-toy/main.go

$(CLIENT):
	go build -v \
	  -tags release \
		-ldflags '-X main.Version=$(VERSION)' \
		-o $(CLIENT) \
  	cmd/web-toy-client/main.go

clean:
	rm -rfv bin

fmt:
	gofmt -l -s -w $(FILES)

test:
	TEST=1 go test $(PACKAGES)
