build:
	go build -C ./src -o ../bin/go_api

run: build
	./bin/go_api

test:
	go test -v ./...