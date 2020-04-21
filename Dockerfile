FROM golang:latest

RUN go get -u github.com/oxequa/realize && \
    go get -u github.com/gin-gonic/gin

WORKDIR /go/src/github.com/home_streaming/

