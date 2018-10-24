FROM golang:alpine as builder

RUN apk add --update git perl-utils bash

WORKDIR $GOPATH/src/github.com/thoughtworks/talisman

RUN go get github.com/mitchellh/gox

VOLUME [$GOPATH/src/github.com/thoughtworks/talisman]

CMD ["/bin/bash", "./build"]


