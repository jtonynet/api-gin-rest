#!/bin/sh

go install github.com/swaggo/swag/cmd/swag@latest
swag init --parseDependency --parseInternal

go run main.go -b 0.0.0.0