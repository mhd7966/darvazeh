FROM golang:1.16-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

ARG BUILD_PATH=./cmd/api

COPY . .

RUN go build -ldflags="-w -s" -o go-app $BUILD_PATH

FROM alpine:3.14.0

COPY --from=builder /app/go-app /app/go-app

CMD ["/app/go-app"]