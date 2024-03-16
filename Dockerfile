FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY internal/app/usecase/testdata/orders_test.json schema.json ./
COPY cmd/app ./cmd/app/
COPY internal ./internal/

RUN go build -o main ./cmd/app

EXPOSE 8080

CMD ["./main"]
