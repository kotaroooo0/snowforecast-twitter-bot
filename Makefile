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

## Watch & run tests
.PHONY: watchtest
watchtest:
	Watch -t make test | cgt

## Build binary
.PHONY: build
build: deps
	go build -o main

## Run binary
.PHONY: run
run: build
	./main

## Watch & run app
.PHONY: realize
realize:
	realize start

## Clean binary
.PHONY: clean
clean:
	go clean
	rm -f main

.PHONY: local
local:
	docker-compose -f docker-compose/docker-compose.local.yml up -d
	sh docker-compose/elasticsearch/setup.sh

## docker-compose up for production
.PHONY: production
production:
	docker-compose -f docker-compose/docker-compose.prod.yml up -d
	sh docker-compose/elasticsearch/setup.sh

## Show help
.PHONY: help
help:
	@make2help $(MAKEFILE_LIST)
