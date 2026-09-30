# =========================
# Build stage
# =========================
FROM golang:1.27.1-alpine AS builder

WORKDIR /src

COPY go.mod .
COPY main.go .

ARG VERSION=dev

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o /out/app .


# =========================
# Runtime stage
# =========================
FROM alpine:3.22

COPY --from=builder /out/app /app

RUN chmod +x /app

EXPOSE 8080

ENTRYPOINT ["/app"]