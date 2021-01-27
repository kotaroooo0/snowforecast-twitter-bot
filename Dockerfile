FROM golang:1.14-alpine as builder

WORKDIR /go/src/github.com/kotaroooo0/snowforecast-twitter-bot/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main

FROM alpine:latest

COPY .env /
COPY batch.snow_resorts.yaml /
COPY wait-for-it.sh /
RUN chmod 777 ./wait-for-it.sh

COPY --from=builder /go/src/github.com/kotaroooo0/snowforecast-twitter-bot/app/main /
EXPOSE 3000

CMD ["./main"]
