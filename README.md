# pg-wait

### Install:
```bash
go install github.com/partyzanex/pg-wait@v0.1.0
```

### Usage:
```bash
$ pg-wait -h
Usage of ~/go/bin/pg-wait:
  -d, --dsn string         postgres URL
  -t, --timeout duration   timeout (default 10s)
  -v, --verbose            verbose errors
pflag: help requested

$ pg-wait --dsn="postgresql://postgres:postgres@localhost:5432/mydatabase?sslmode=disable" -t 5s -v 
2021/04/11 00:10:18 try tcp connection: dial tcp [::1]:5432: connectex: No connection could be made because the target machine actively refused it.
2021/04/11 00:10:21 try tcp connection: dial tcp [::1]:5432: connectex: No connection could be made because the target machine actively refused it.
2021/04/11 00:10:26 exit by timeout
exit status 1

$ pg-wait --dsn="postgresql://postgres:postgres@localhost:5432/mydatabase?sslmode=disable" -t 5s -v
2021/04/11 00:13:52 success!
```

### Usage via docker-compose:

[Dockerfile](https://github.com/partyzanex/pg-wait/blob/main/example/Dockerfile):
```Dockerfile
FROM golang:1.16-alpine as builder

WORKDIR /go/src/pg-wait

ENV GOPATH /go
ENV PG_WAIT_VERSION v0.1.0

RUN go mod init fake && go install github.com/partyzanex/pg-wait@${PG_WAIT_VERSION}


FROM alpine:3.13

COPY --from=builder /go/bin/pg-wait /usr/local/bin
```

[docker-compose.yml](https://github.com/partyzanex/pg-wait/blob/main/example/docker-compose.yml):
```yaml
version: '3.5'

services:
  postgres:
    image: postgres:12-alpine
    environment:
      POSTGRES_USER: "test"
      POSTGRES_PASSWORD: "test"
      POSTGRES_DB: "test"
    networks:
      - example

  pg-wait:
    build:
      dockerfile: Dockerfile
      context: .
    command: sh -c "pg-wait --dsn=\"postgresql://test:test@postgres:5432/test?sslmode=disable\" --verbose --timeout=10s"
    depends_on:
      - postgres
    networks:
      - example

networks:
  example:
    name: example_net
```

Run:
```bash
$ docker-compose up --build
...
Creating example_postgres_1 ... done
Creating example_pg-wait_1  ... done
Attaching to example_postgres_1, example_pg-wait_1
pg-wait_1   | 2021/04/10 21:48:16 try tcp connection: dial tcp 172.22.0.2:5432: connect: connection refused
postgres_1  | The files belonging to this database system will be owned by user "postgres".
...
postgres_1  | 2021-04-10 21:48:17.772 UTC [1] LOG:  database system is ready to accept connections
pg-wait_1   | 2021/04/10 21:48:18 success!
example_pg-wait_1 exited with code 0
```