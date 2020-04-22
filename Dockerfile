FROM golang:1.14.2-alpine3.11 as build

LABEL maintainer="https://gihub.com/popper2710"

RUN set -ex && \
    apk update && \
    apk add --no-cache git && \
    apk add gcc libc-dev && \
    go get -u github.com/gin-gonic/gin && \
    go get -u github.com/oxequa/realize

WORKDIR /go/src/github.com/home_streaming

FROM alpine:latest

WORKDIR /app

COPY --from=build /go/src/github.com/home_streaming/app .

RUN set -x && \
    addgroup go && \
    adduser -D -G go go && \
    chown -R go:go /app/app
