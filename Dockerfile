FROM golang:1.21.0-alpine as builder

WORKDIR /app

COPY ./ ./

RUN go build -o ./main ./cmd/namer/main.go

RUN rm -rfv !("main")

CMD ["/app/main"]
