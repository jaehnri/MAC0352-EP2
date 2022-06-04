FROM golang:1.17-alpine as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o client ./cmd/client

EXPOSE 8080

#CMD ["./client 172.17.0.3 8080 tcp"]
