FROM golang:1.17-alpine as builder

RUN apk add --no-cache git

WORKDIR /go/src/pg-wait

ENV GOPATH /go
ENV PG_WAIT_VERSION v0.1.1

RUN go mod init fake && go install github.com/partyzanex/pg-wait@${PG_WAIT_VERSION}


FROM alpine:3.13

COPY --from=builder /go/bin/pg-wait /usr/local/bin
