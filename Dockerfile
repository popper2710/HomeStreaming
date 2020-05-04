FROM golang:1.14.2-alpine3.11 as build

LABEL maintainer="https://gihub.com/popper2710"

RUN set -ex && \
    apk update && \
    apk add --no-cache git ffmpeg bash&&\
    go get -u github.com/oxequa/realize

COPY ./app /go/src/github.com/home_streaming/app

WORKDIR /go/src/github.com/home_streaming/app

RUN go mod download

FROM alpine:latest

WORKDIR /app

COPY --from=build /go/src/github.com/home_streaming/app .

RUN set -x && \
    addgroup go && \
    adduser -D -G go go && \
    chown -R go:go /app/app

ENV GIN_MODE release
