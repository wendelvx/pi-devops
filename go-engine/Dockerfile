# Estágio de Compilação
# Ajustado de 1.25 para 1.24-alpine (a versão 1.25 ainda não é estável/disponível)
FROM golang:1.24-alpine AS builder

# Instalação de dependências de build necessárias para alguns pacotes Go
RUN apk add --no-cache git

WORKDIR /src

# Copiamos apenas os arquivos de dependência primeiro para otimizar o cache do Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiamos o restante do código fonte
COPY . .

# Compilação do binário estático (CGO_ENABLED=0 garante que ele não dependa de bibliotecas C do host)
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/go-engine ./cmd/engine

# Estágio Final (Runtime)
FROM alpine:3.21

# Instalação de certificados CA (essencial para que o Go consiga se comunicar via HTTPS 
# com o Cloudflare Tunnel ou Azure, conforme descrito no seu plano de arquitetura)
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copiamos o binário gerado no estágio anterior
COPY --from=builder /bin/go-engine /usr/local/bin/go-engine

# Exposição da porta configurada no seu mux
EXPOSE 8080

# Execução do binário usando o caminho absoluto para evitar problemas de PATH
CMD ["/usr/local/bin/go-engine"]