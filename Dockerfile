FROM golang:1.16.7-alpine3.13
WORKDIR /go/src
RUN apk --no-cache add make gcc musl-dev

COPY . .
RUN make build
