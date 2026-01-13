# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main cmd/api/main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY start.sh .
COPY internal/db/migrate_files ./internal/db/migrate_files

# Install migrate tool
RUN apk add --no-cache curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate /usr/bin/migrate

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
