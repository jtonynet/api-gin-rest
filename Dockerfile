FROM golang:1.21.1-alpine

WORKDIR /usr/src/app

COPY . . 

RUN go install github.com/swaggo/swag/cmd/swag@latest

CMD ["go", "run", "main.go", "-b", "0.0.0.0"]