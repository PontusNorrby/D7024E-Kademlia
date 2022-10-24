# Peer-to-Peer Distributed Data Store

For this lab we use Goland and Docker compose.
For information about downloading and installing Golang on your computer see [go website](https://go.dev/).
For information about downloading and installing docker compose on your computer see [docker website](https://docs.docker.com/compose/install/).

##### To run the code, follow the following steps:
Head to a folder of your choice on your computer and in the cmd run:
```
git clone https://github.com/PontusNorrby/D7024E-Kademlia.git
```
Now head to the folder you cloned the code in and in the cmd run:
```
docker build . -t kadlab
```
And after that:
```
docker-compose up
```
Your containers are now up and running.

##### To run the CLI:
Head to the folder you cloned the code in and run:
```
docker build . -t kadlab
```
And then run:
```
docker-compose up -d
```
Now open docker compose and choose a docker container, copy the ID of the container and run:
```
docker attach [the ID goes here ]
```
Now you have the CLI up. To get a list of available commands write help in the cmd

##### To run the testes:
Head to the folder you cloned the code in and the to src folder and then to kademlia folder and run:
```
go test
```
