FROM golang:alpine3.21 AS builder

WORKDIR /usr/app

COPY ./ ./

RUN apk -U upgrade
RUN apk add gcc make libc-dev

ENV CGO_ENABLED=1

RUN make build/server
RUN make build/cli

FROM alpine:3.21

ENV CACHER_PERSISTANCE=false

COPY --from=builder /usr/app/bin /usr/local/bin/cacher

RUN apk -U upgrade
RUN apk add netcat-openbsd

EXPOSE 8595

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 CMD \
  echo "GET foo" | nc localhost 8595 || exit 1

CMD [ "sh", "-c", "/usr/local/bin/cacher/server -persist=${CACHER_PERSISTANCE}"]
