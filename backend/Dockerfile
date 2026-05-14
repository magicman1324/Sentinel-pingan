FROM golang:1.25-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /monitor-backend ./cmd/server/

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /monitor-backend /monitor-backend
EXPOSE 8080 9090
ENTRYPOINT ["/monitor-backend"]
