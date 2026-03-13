FROM mirror.gcr.io/library/golang:1.25-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git ca-certificates

ARG GOPROXY=https://proxy.golang.org,direct
ENV GOPROXY=${GOPROXY}

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o engine ./cmd/engine

FROM mirror.gcr.io/library/alpine:3.19
WORKDIR /app
COPY --from=builder /app/engine ./engine

EXPOSE 8080
ENTRYPOINT ["./engine"]
