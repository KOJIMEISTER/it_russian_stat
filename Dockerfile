FROM golang:1.23.4-bookworm AS deps
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM deps AS api-builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/api ./cmd/api

FROM deps AS scrapper-builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/scrapper ./cmd/scrapper

FROM gcr.io/distroless/static-debian12 as api-final
WORKDIR /
COPY --from=api-builder --chown=nonroot:nonroot /app/bin/api /app/api
USER nonroot
CMD [ "/app/api" ]

FROM gcr.io/distroless/static-debian12 as scrapper-final
WORKDIR /
COPY --from=scrapper-builder --chown=nonroot:nonroot /app/bin/scrapper /app/scrapper
USER nonroot
CMD [ "/app/scrapper" ]