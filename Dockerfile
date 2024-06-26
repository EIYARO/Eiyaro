# Build Eiyaro in a stock Go builder container
FROM golang:1.9-alpine as builder

RUN apk add --no-cache make git

ADD . /go/src/github.com/eiyaro
RUN cd /go/src/github.com/eiyaro && make eiyarod && make eiyarocli

# Pull Eiyaro into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/eiyaro/cmd/eiyarod/eiyarod /usr/local/bin/
COPY --from=builder /go/src/eiyaro/cmd/eiyarocli/eiyarocli /usr/local/bin/

EXPOSE 1999 46656 46657 9888
