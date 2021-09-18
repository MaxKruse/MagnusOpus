FROM golang:1.16-alpine

# Install dependencies
RUN apk add build-base

WORKDIR /app

RUN go get github.com/cosmtrek/air

COPY go.mod ./
COPY go.sum ./

RUN go mod download
