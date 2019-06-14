FROM golang:alpine

MAINTAINER 黄滚<hg@yuanli-inc.com> 

RUN apk add --no-cache tini mariadb-client

ADD . /go/src/github.com/yuanfenxi/yuanlicast

RUN apk add --no-cache mariadb-client
RUN cd /go/src/github.com/yuanfenxi/yuanlicast/ && \
    go build -o bin/yuancast ./cmd/go-mysql-elasticsearch && \
    cp -f ./bin/yuancast /go/bin/yuancast

ENTRYPOINT ["/sbin/tini","--","yuancast"]
