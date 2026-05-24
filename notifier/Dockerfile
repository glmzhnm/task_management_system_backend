FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod init notifier && go build -o notifier-app .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/notifier-app .
EXPOSE 8081
CMD ["./notifier-app"]
