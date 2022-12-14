# syntax=docker/dockerfile:1

FROM golang:latest


WORKDIR $GOPATH/src/github.com/PontusNorrby/D7024E-Kademlia

RUN apt update
RUN apt install reptyr
RUN yes yes | apt install screen

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build src/main.go

CMD  ["./main"]