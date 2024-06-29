build:
	go build -o ./bin/boriqua_project

run: build
	./bin/boriqua_project

test:
	go test ./...