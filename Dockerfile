FROM golang:1.23.4-bookworm AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/main ./cmd/api

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=builder --chown=nonroot:nonroot /app/bin/main /app/main
USER nonroot
EXPOSE 8080
ENTRYPOINT [ "/app/main" ]