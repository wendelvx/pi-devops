FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/go-engine ./cmd/engine

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /bin/go-engine /usr/local/bin/go-engine

EXPOSE 8080

CMD ["go-engine"]
