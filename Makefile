GOPATH := B:/Projects/Go
GOBIN := $(GOPATH)/bin
export GOPATH
export GOBIN

build:
	go build -o $(GOBIN)/blufi-network ./cmd/node

run: build
	$(GOBIN)/blufi-network
