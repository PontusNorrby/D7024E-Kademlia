# syntax=docker/dockerfile:1

FROM golang:latest


WORKDIR $GOPATH/src/github.com/PontusNorrby/D7024E-Kademlia

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build src/main.go
EXPOSE 3000

CMD  ["./main"]
