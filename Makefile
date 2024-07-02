GOPATH := B:\Projects\Go
GOBIN := $(GOPATH)/bin
export GOPATH
export GOBIN

build:
	cd $(GOPATH) && go build -o $(GOBIN)/boriqua_project

run: build
	$(GOBIN)/boriqua_project

test:
	go test ./...