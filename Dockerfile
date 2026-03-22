FROM golang:1.26.1-alpine3.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /roomly ./cmd/roomly/main.go

FROM alpine:latest
COPY --from=builder /roomly /roomly
ENTRYPOINT [ "/roomly" ]
