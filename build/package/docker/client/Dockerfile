FROM golang:1.17-alpine as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o client ./cmd/client

EXPOSE 8080
