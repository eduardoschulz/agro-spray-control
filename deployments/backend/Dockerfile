FROM golang:1.24 AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./backend/
WORKDIR /app/backend
RUN go mod download

COPY backend/ ./

RUN go build -o /app/main ./cmd 


FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 30000
CMD ["./main"]
