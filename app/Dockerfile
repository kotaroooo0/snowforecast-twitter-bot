FROM golang as builder

WORKDIR /go/src/github.com/kotaroooo0/snowfall-bot/app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o main

FROM scratch

COPY --from=builder /go/src/github.com/kotaroooo0/snowfall-bot/app/main /main
COPY .env /

EXPOSE 3000

ENTRYPOINT ["/main"]
