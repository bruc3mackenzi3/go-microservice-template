FROM golang:1.18-alpine

EXPOSE 80
WORKDIR /app

COPY ./ ./

RUN go mod download
RUN go build -o microservice-app

CMD [ "./microservice-app" ]
