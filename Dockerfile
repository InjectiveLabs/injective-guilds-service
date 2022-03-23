FROM golang:1.17.3-alpine as builder
RUN apk add --no-cache git make gcc

WORKDIR /src
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /bin/injective-guilds ./cmd/injective-guilds/
EXPOSE 9920
CMD ["injective-guilds"]
