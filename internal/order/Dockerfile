FROM golang:1.24-alpine

WORKDIR /order_service

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o order-service .

CMD ["./order-service"]