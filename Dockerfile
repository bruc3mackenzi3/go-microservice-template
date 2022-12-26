FROM golang:1.18-alpine

WORKDIR /app

COPY ./users-docker ./

CMD [ "./users-docker" ]
