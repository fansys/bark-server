FROM golang:1.16-alpine AS builder

ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn

COPY . /go/src/bark-server

WORKDIR /go/src/bark-server

RUN set -ex \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && apk add git gcc libc-dev \
    && BUILD_VERSION=$(cat version) \
    && BUILD_DATE=$(date "+%F %T") \
    && COMMIT_SHA1=$(git rev-parse HEAD) \
    && go install -trimpath -ldflags \
            "-X 'main.version=${BUILD_VERSION}' \
             -X 'main.buildDate=${BUILD_DATE}' \
             -X 'main.commitID=${COMMIT_SHA1}' \
             -w -s"

FROM alpine:3.13

LABEL maintainer="mritd <mritd@linux.com>"

# set up nsswitch.conf for Go's "netgo" implementation
# - https://github.com/golang/go/blob/go1.9.1/src/net/conf.go#L194-L275
# - docker run --rm debian:stretch grep '^hosts:' /etc/nsswitch.conf
RUN [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf

RUN set -ex \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && apk upgrade \
    && apk add ca-certificates

COPY --from=builder /go/bin/bark-server /usr/local/bin/bark-server
COPY deploy/AuthKey.p8 /usr/local/bin/

VOLUME /data

EXPOSE 8080

CMD ["bark-server"]
