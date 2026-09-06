FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN go build -o stockstalk ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/stockstalk .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/mock ./mock
EXPOSE 8080
CMD ["./stockstalk"]
