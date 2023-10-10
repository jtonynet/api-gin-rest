FROM golang:1.21.1-alpine

WORKDIR /usr/src/app

COPY . . 

RUN go install github.com/swaggo/swag/cmd/swag@latest
