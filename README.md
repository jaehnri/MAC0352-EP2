# Tic Tac Toe - MAC0352-EP2
Tic tac toe game over the internet.

Second project for MAC0352 - Computer Networks and Distributed Systems.

## Build and Run
You need golang installed to run this program.

**Run Client:**
```sh
make build-server
make run-server
```

**Run Server:**
```sh
make build-client
make run-client
```

## Docker Compose:

The easiest way to run the server and the database is to run:
```sh
docker-compose up -d
```

It makes sure to run the server, the database and the database volume so the data is kept even if the database is killed.

It deploys the containers on Docker's bridge network. It assigns IP addresses to the containers. It can be inspected like this:
```sh
docker inspect network bridge
```

### Cleaning

To kill the server and the database:
```sh
docker-compose down
```

To kill everything (including the database volume):
```sh
docker-compose down --remove-orphans --volumes
```

## Test
```sh
make test
```
