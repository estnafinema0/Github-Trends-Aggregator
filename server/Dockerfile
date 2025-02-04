FROM golang:1.23-alpine

LABEL maintainer="estnafinema0@gmail.com"

WORKDIR /app

RUN apk update && apk upgrade && apk add bash

COPY . .

EXPOSE 8080

CMD ["go", "run", "main.go"]