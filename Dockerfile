FROM golang:alpine as builder

RUN apk add --update git perl-utils ca-certificates && \
    update-ca-certificates && \
	mkdir -p /talisman-src && \
	wget https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz && \
	mkdir -p /opt && \
	tar -xf upx-3.96-amd64_linux.tar.xz -C /opt && \
	git config --global user.name "Talisman Maintainers" && \
	git config --global user.email "talisman-maintainers@thoughtworks.com "

ENV CGO_ENABLED=0
ENV PATH="$PATH:/opt/upx-3.96-amd64_linux"

WORKDIR /talisman-src

CMD ["build-release-binaries"]
