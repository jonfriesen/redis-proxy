FROM golang:alpine as builder

LABEL maintainer "Jon Friesen <jon@jonfriesen.ca>"

RUN apk add build-base gcc abuild binutils binutils-doc gcc-doc git redis

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN mkdir -p /go/src/github.com/jonfriesen/redis-proxy
ADD . /go/src/github.com/jonfriesen/redis-proxy

WORKDIR /go/src/github.com/jonfriesen/redis-proxy

RUN chmod +x scripts/wait-for-redis.sh

RUN go get -v ./...
RUN go build ./cmd/redis-proxy-http

FROM alpine

RUN apk add redis

COPY --from=builder /go/src/github.com/jonfriesen/redis-proxy/redis-proxy-http ./usr/local/bin/redis-proxy-http
COPY --from=builder /go/src/github.com/jonfriesen/redis-proxy/scripts/wait-for-redis.sh ./wait-for-redis.sh
ENTRYPOINT ["./wait-for-redis.sh", "redis-proxy-http", "-redis-host=redis"]