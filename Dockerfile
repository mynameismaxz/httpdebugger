# Stage 1: Build the application
FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o server cmd/server/main.go

# Stage 2: Create a minimal image to run the application
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 1337

CMD ["./server"]
