FROM golang:1.24-alpine3.20 AS build

WORKDIR /opt/app

RUN apk add curl

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN mkdir bin

RUN go build -ldflags '-s' -o bin ./cmd/...

FROM alpine:3.20

WORKDIR /opt/app

COPY --from=build /opt/app/bin ./bin

