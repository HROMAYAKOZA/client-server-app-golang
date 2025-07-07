FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/client/main.go ./
COPY internal/client/ internal/client/
RUN go build -o client ./main.go

# чистим кеш, уменьшаем размер образа
RUN go clean -modcache

CMD ["./client"]
