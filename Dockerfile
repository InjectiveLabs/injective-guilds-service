#install packages for build layer
FROM golang:1.17.3-alpine as builder
RUN apk add --no-cache git make gcc libc-dev linux-headers

#build binary
WORKDIR /src
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make install

#build main container
FROM alpine:latest
RUN apk add --update --no-cache ca-certificates
RUN apk add curl
COPY --from=builder /go/bin/* /usr/local/bin/
EXPOSE 9930
