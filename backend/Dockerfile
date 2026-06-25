# ── base ─────────────────────────────────────────────────────────────────────
FROM golang:1.26.4-alpine AS base
WORKDIR /app


# ── mod-download ──────────────────────────────────────────────────────────────
# Isolated dependency layer — only invalidated when go.mod or go.sum change,
# not on every source change.
FROM base AS mod-download
COPY go.mod go.sum ./
RUN go mod download


# ── builder ───────────────────────────────────────────────────────────────────
# Compiles the binary. Inherits the warm module cache from mod-download.
FROM mod-download AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/glimpse


# ── final ─────────────────────────────────────────────────────────────────────
# Minimal runtime image — no Go toolchain, no source, no module cache.
# Only the compiled binary and ca-certificates for outbound TLS.
FROM alpine:3.21 AS final
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/main .

CMD ["./main"]