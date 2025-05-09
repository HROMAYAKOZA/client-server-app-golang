FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY client/client.go ./
RUN go build -o client ./client.go

# чистим кеш, уменьшаем размер образа
RUN go clean -modcache

CMD ["./client"]
