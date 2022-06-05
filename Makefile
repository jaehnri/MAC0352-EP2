########################################################################################################################
# Used to compile the server
build-server:
	go build -o server ./cmd/server

# Used to run the server locally
run-server:
	./server

########################################################################################################################
# Used to compile the client
build-client:
	go build -o client ./cmd/client

# Used to run the client locally
run-client:
	./client localhost 8080 tcp

########################################################################################################################
# Execute unit tests
test:
	go test ./...

########################################################################################################################
# Deploy the database and the server on the docker network
run:
	docker-compose up -d --build

# Kill the containers and remove the local images
shutdown:
	docker-compose down --remove-orphans --rmi local

# Kill the containers, rebuild the image and deploy
restart: shutdown run

# Kill the containers and remove the volume (database data)
remove-all:
	docker-compose down --remove-orphans --rmi local --volumes

########################################################################################################################
# Build the client Docker image
build-client-docker:
	docker build -f build/package/docker/client/Dockerfile . -t client-docker

# Run the dockerized client using TCP
run-client-docker-tcp-interactive: build-client-docker
	docker run -i client-docker ./client 172.17.0.3 8080

# Run the dockerized client using UDP
run-client-docker-udp-interactive: build-client-docker
	docker run -i client-docker ./client 172.17.0.3 8080 udp

########################################################################################################################
# Run the dockerized client using TCP
run-client-docker-tcp: build-client-docker
	docker run client-docker

########################################################################################################################
# Test with 0 clients
test-0-clients:
	rm -rf test-outputs/0/
	mkdir test-outputs/0/
	./scripts/tests/0/test-0-clients.sh test-outputs/0/cpu-0-test test-outputs/0/network-0-test

test-2-clients:
	rm -rf test-outputs/2/
	mkdir test-outputs/2/
	./scripts/tests/2/test-2-active-clients.sh test-outputs/2/cpu-2-test test-outputs/2/network-2-test 1
