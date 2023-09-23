FROM golang:1.21.1

WORKDIR /usr/src/app

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN swag init --parseDependency --parseInternal

