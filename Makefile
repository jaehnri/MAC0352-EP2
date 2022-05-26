build-server:
	go build -o server ./cmd/server

run-server:
	./server

build-client:
	go build -o client ./cmd/client

run-client:
	./client

test:
	go test ./...
