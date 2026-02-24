# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/smarterp ./cmd/api

FROM alpine:3.20
RUN adduser -D -u 10001 appuser
USER appuser

WORKDIR /app
COPY --from=build /bin/smarterp /app/smarterp

EXPOSE 8080
ENTRYPOINT ["/app/smarterp"]
