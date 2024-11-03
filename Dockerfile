FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o mustika .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/mustika .

CMD ["./mustika"]