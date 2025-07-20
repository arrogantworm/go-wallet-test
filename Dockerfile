FROM golang:1.24.5-bookworm

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o wallet-backend ./cmd/main.go

CMD ["/wallet-backend"]
