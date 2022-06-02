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

run:
	docker-compose up -d --build

restart: shutdown run

shutdown:
	docker-compose down --remove-orphans --rmi local

remove-all:
	docker-compose down --remove-orphans --rmi local --volumes