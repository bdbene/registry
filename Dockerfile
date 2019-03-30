FROM golang:1.9.1 AS builder

WORKDIR $GOPATH/src/github.com/bdbene/registry
COPY . .

RUN go get -d -v ./...
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/registry
COPY config.tml $GOPATH/bin/

FROM debian:jessie-20190326 AS package

COPY --from=builder /go/bin/registry/* /root/
RUN mkdir /data

ENTRYPOINT ["./registry"]