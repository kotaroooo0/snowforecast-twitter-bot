FROM golang:1.14-alpine as builder

WORKDIR /go/src/github.com/kotaroooo0/snowforecast-twitter-bot/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main

FROM alpine:latest

COPY --from=builder /go/src/github.com/kotaroooo0/snowforecast-twitter-bot/app/main /main
COPY .env /
COPY batch.snow_resorts.yaml /

EXPOSE 3000

ENTRYPOINT ["/main"]
