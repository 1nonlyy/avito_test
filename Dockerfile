FROM golang:1.24.2 

WORKDIR /app

COPY . .

RUN go mod tidy

# Добавим вывод списка файлов — DEBUG
RUN ls -la ./cmd/server

# Сборка
RUN GOARCH=arm64 GOOS=linux go build -o main ./cmd/server

# Проверим результат сборки


EXPOSE 8080
CMD ["/app/main"]

