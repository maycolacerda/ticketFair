
FROM golang:1.22-alpine as builder


RUN apk update && apk upgrade --no-cache

WORKDIR /app
# Copia os arquivos go.mod e go.sum para o diretório de trabalho
COPY go.mod go.sum ./
RUN go mod download

#instala o curl para o health check
RUN apk add --no-cache curl
# Copia o restante dos arquivos do projeto para o diretório de trabalho
COPY . .

# Compila o aplicativo Go e gera um executável chamado ticketfair
RUN go build -o ticketfair -ldflags "-s -w" .

FROM alpine:latest

# Define o diretório de trabalho
WORKDIR /app

# Copia o executável compilado do estágio 'builder'
COPY --from=builder /app/ticketfair .

# expoe a porta 8000 para acesso externo
EXPOSE 8000

# Define o comando padrão a ser executado quando o container iniciar
CMD ["./ticketfair"]
