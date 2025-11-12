FROM golang:1.25 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /usr/local/bin/app ./cmd/web/

FROM debian:bookworm-slim

WORKDIR /usr/local/bin/
COPY ./ui /usr/local/bin/ui
COPY --from=builder /usr/local/bin/app /usr/local/bin/app


CMD ["app"]
