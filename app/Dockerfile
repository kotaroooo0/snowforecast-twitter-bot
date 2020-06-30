# syntax = docker/dockerfile:experimental
FROM golang as builder

WORKDIR /go/src/github.com/kotaroooo0/snowforecast-twitter-bot/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o main

FROM scratch

COPY --from=builder /go/src/github.com/kotaroooo0/snowforecast-twitter-bot/app/main /main
COPY .env /

EXPOSE 3000

ENTRYPOINT ["/main"]
