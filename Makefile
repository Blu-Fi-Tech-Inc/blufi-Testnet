GOPATH := B:/Projects/Go
GOBIN := $(GOPATH)/bin
GOCACHE := $(GOPATH)/cache
export GOPATH
export GOBIN
export GOCACHE

build:
	go build -o $(GOBIN)/boriqua_project ./cmd/node

run: build
	$(GOBIN)/boriqua_project
