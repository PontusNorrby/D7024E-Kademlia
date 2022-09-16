Course D7024E At LTU

Laboration in Creating a Peer-to-Peer Distributed Data Store using the Kademlia algorithm.

# Docker File
---
We got the Docker file from the official Docker [documentation](https://docs.docker.com/language/golang/build-images/)
```
FROM alpine:latest
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./  
RUN go build -o /kademlia 
EXPOSE 8080
CMD [ "/kademlia" ]
```

### Build the image
To build the image, head to the folder you have your docker file and use the following command:

```
docker build --tag kademlia .
```

---