FROM golang:alpine3.21 AS builder

WORKDIR /usr/app

COPY ./ ./

RUN apk -U upgrade

RUN go build -o ./bin/server -v -race ./cmd/server
RUN go build -o ./bin/cli -v -race ./cmd/cli

FROM alpine:3.21

COPY --from=builder /usr/app/bin /usr/local/bin/cacher

RUN apk -U upgrade
RUN apk add netcat-openbsd

EXPOSE 8595

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 CMD \
  echo "GET foo" | nc localhost 8595 || exit 1

CMD [ "/usr/local/bin/cacher/server" ]
