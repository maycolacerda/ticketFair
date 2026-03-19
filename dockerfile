FROM golang:1.26-alpine3.22 AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    $(go env GOPATH)/bin/swag init -g main.go

RUN CGO_ENABLED=0 GOOS=linux GOFLAGS="-trimpath" go build -ldflags="-w -s" -o ticketfair .

FROM alpine:3.21

WORKDIR /app

RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

COPY --from=builder --chown=appuser:appuser /app/ticketfair .
COPY --from=builder --chown=appuser:appuser /app/docs ./docs

USER appuser

EXPOSE 8000

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8000/api/v1 || exit 1

CMD ["./ticketfair"]
