FROM golang:1.18-alpine AS builder

WORKDIR /
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN go build

FROM alpine:3.15
WORKDIR /
COPY --from=builder /device-plugin-1 .
ENTRYPOINT ["/device-plugin-1"]
