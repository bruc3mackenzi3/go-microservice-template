FROM gcr.io/distroless/base-debian10

WORKDIR /app

COPY ./users-docker ./

CMD [ "./users-docker" ]
