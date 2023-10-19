FROM golang:1.20 as builder

WORKDIR /usr/src/app

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
