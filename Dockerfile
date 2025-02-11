# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod .
# go.sum がない場合でもエラーにならないように条件付きでコピー
COPY go.sum* .
RUN go mod download
COPY . .
RUN go build -o todoapp ./main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/todoapp .
EXPOSE 8080
CMD ["./todoapp"]