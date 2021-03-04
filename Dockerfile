FROM golang:alpine as builder

RUN apk add --update git perl-utils bash

RUN mkdir -p /talisman-src

RUN go get github.com/mitchellh/gox && wget https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz && mkdir -p /opt && tar -xf upx-3.96-amd64_linux.tar.xz -C /opt

COPY ./build-release-binaries /usr/local/bin/

ENV PATH="$PATH:/opt/upx-3.96-amd64_linux"

WORKDIR /talisman-src

CMD ["build-release-binaries"]
