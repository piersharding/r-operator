FROM golang:1.12 AS build

COPY . /go/src/r-controller
WORKDIR /go/src/r-controller
RUN go get -u github.com/golang/dep/cmd/dep && dep ensure && go build -o /go/bin/r-controller

FROM debian:stretch-slim

COPY --from=build /go/bin/r-controller /usr/bin/r-controller
