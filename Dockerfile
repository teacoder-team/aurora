# Этап сборки
FROM golang:1.23-alpine as builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для зависимостей
COPY go.mod ./
COPY go.sum ./

# Загружаем зависимости
RUN go mod tidy

# Копируем исходный код приложения
COPY . .

# Собираем бинарник, указывая путь к главному файлу
RUN go build -o main ./cmd/main.go

# Этап выполнения
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем бинарник из этапа сборки
COPY --from=builder /app/main .

# Копируем .env файл для конфигурации среды
COPY .env .env

# Копируем config.toml в контейнер
COPY config.toml ./config.toml

# Открываем порт, на котором будет работать приложение
EXPOSE 14704

# Запуск приложения
CMD ["./main"]
