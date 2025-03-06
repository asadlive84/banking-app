
FROM golang:1.23-alpine as builder

WORKDIR /app


COPY go.mod go.sum ./

RUN go mod download


COPY . .


RUN go build -o banking_app .

FROM golang:1.23-alpine

WORKDIR /app

COPY --from=builder /app/banking_app .
COPY --from=builder /app/config ./config

CMD ["./banking_app"]
