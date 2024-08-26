# Stage 1: Build the Go binary
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -v -o bin/mup .

# Stage 2: Create a lightweight image for running the binary
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/mup ./bin/mup

EXPOSE 5000

CMD ["./bin/mup"]
