FROM golang:1.17-alpine

# Example abstracted from https://docs.docker.com/language/golang/build-images/

WORKDIR /app

COPY ./ ./

RUN go mod download
RUN go build -o microservice-app

EXPOSE 8080

CMD [ "./microservice-app" ]
