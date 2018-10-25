FROM golang:alpine as e2eTester

LABEL maintainer "Jon Friesen <jon@jonfriesen.ca>"

RUN apk add build-base gcc abuild binutils binutils-doc gcc-doc git

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN mkdir -p /go/src/github.com/jonfriesen/redis-proxy
ADD . /go/src/github.com/jonfriesen/redis-proxy

WORKDIR /go/src/github.com/jonfriesen/redis-proxy

RUN go get -v ./...
# RUN go test -tags=e2e ./...
ENTRYPOINT [ "go", "test", "-tags=e2e", "./..." ]