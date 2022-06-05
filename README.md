# Tic Tac Toe - MAC0352-EP2
Tic tac toe game over the internet.

Second project for MAC0352 - Computer Networks and Distributed Systems.

# Build and Run

## Docker Compose:

The necessary dependencies required to run this project are:
- [docker-compose](https://docs.docker.com/compose/install/)
- [docker](https://www.docker.com/)
- [golang](https://go.dev/dl/)

Also, it is required that your host is able to access the `docker network bridge` as the database and the server are deployed there.
This network is useful so we can assign IP addresses to the containers. It can be inspected like this:
```sh
docker inspect network bridge
```

If you have `docker` installed, and you are using a Linux environment, it is almost certain that this network is running and available.

The easiest way to run the server and the database is to run:
```sh
make run
```

It makes sure to run the server, the database and the database volume so the data is kept even if the database is killed.

Notice that this system assumes that the IPs `172.17.0.2` and `172.17.0.3` are available for use of the database and the server respectively.
In case they aren't, the server won't be able to run.

After running, you can check if everything is ok by checking if the containers are up:
```sh
docker-compose ps
```

And also by checking the server logs:
```sh
docker logs -f mac0352-ep2_server_1
```

### Observations about Mac

There is no `docker0` bridge on macOS. Because of the way networking is implemented in Docker Desktop for Mac, you cannot see a `docker0` interface on the host. This interface is actually within the virtual machine.

Thus, it is a known Docker for Mac issue that the [the docker (Linux) bridge network is not reachable from the macOS host](https://docs.docker.com/desktop/mac/networking/#per-container-ip-addressing-is-not-possible) and [Docker Desktop for Mac can’t route traffic to containers](https://docs.docker.com/desktop/mac/networking/#i-cannot-ping-my-containers).

Basically, this means that this EP won't work out of the box on MacOS environments.

## Running clients
To run a client, build it and run the following:
```sh
./client <server-ip> <server-port> <conn-type>
# <conn-type>: tcp or udp
```
Note that you cannot run multiple client application in the same machine. This is because it listens to the oponent connection in a constant port.

## Cleaning

To kill the server and the database:
```sh
make shutdown
```

To kill everything (including the database volume):
```sh
docker-compose down --remove-orphans --volumes
```


## Test
```sh
make test
```
