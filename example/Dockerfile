FROM golang:1.16-alpine as builder

WORKDIR /go/src/pg-wait

ENV GOPATH /go
ENV PG_WAIT_VERSION v0.1.0

RUN go mod init fake && go install github.com/partyzanex/pg-wait@${PG_WAIT_VERSION}


FROM alpine:3.13

COPY --from=builder /go/bin/pg-wait /usr/local/bin
