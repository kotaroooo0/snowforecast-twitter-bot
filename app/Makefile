NAME := snowforecast-twitter-bot
VERSION := $(gobump show -r)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(REVISION)'

## Install dependencies
.PHONY: deps
deps:
	go get -v -d

## Run tests
.PHONY: test
test: deps
	go test ./... -v

## Show help
.PHONY: help
help:
	@make2help $(MAKEFILE_LIST)

## Build binary
.PHONY: build
build: deps
	go build -o main

## Run binary
.PHONY: run
run: build
	./main

## Clean binary
.PHONY: clean
clean:
	go clean
	rm -f main
