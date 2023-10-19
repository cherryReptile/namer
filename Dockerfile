FROM golang:1.20 as builder

WORKDIR /app

COPY ./ ./

RUN go build -o ./main ./cmd/users/main.go

CMD ["/app/main"]
