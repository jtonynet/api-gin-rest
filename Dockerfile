FROM golang:1.21.1-alpine AS api
WORKDIR /usr/src/app
COPY . . 
RUN go install github.com/swaggo/swag/cmd/swag@latest

FROM golang:1.21.1-alpine AS worker
WORKDIR /usr/src/app
COPY . . 
