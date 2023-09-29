FROM golang:latest

WORKDIR /usr/src/app

COPY . .

# RUN go install github.com/swaggo/swag/cmd/swag@latest

# RUN swag init --parseDependency --parseInternal

