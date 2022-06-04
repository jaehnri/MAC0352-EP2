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

build-client-docker:
	docker build -f build/package/docker/client/Dockerfile . -t client-docker

run-client-docker: build-client-docker
	docker run -i client-docker ./client 172.17.0.3 8080