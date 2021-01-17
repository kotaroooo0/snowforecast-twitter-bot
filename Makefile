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
	redis-cli -n 1 flushdb
	cat db/data.txt | redis-cli -n 1 --pipe
	REDIS_HOST=localhost go test ./... -v

## Watch & run tests
.PHONY: watchtest
watchtest:
	REDIS_HOST=localhost Watch -t make test | cgt

## Build binary
.PHONY: build
build: deps
	go build -o main

## Run binary
.PHONY: run
run: build
	REDIS_HOST=localhost ./main

## Watch & run app
.PHONY: realize
realize:
	realize start

## Init data
.PHONY: initdata
initdata:
	redis-cli flushdb
	cat db/data.txt | redis-cli --pipe

## Run redis
.PHONY: redisup
redisup:
	redis-server /usr/local/etc/redis.conf

## Clean binary
.PHONY: clean
clean:
	go clean
	rm -f main

.PHONY: local
local:
	docker-compose -f docker-compose/docker-compose.local.yml up -d
	sh docker-compose/elasticsearch/setup.sh

## docker-compose up for produciton
.PHONY: produciton
produciton:
	docker-compose -f docker-compose/docker-compose.prod.yml up -d
	sh docker-compose/elasticsearch/setup.sh

## Show help
.PHONY: help
help:
	@make2help $(MAKEFILE_LIST)
