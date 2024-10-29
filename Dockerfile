# Этап сборки (builder stage)
FROM golang:1.23-alpine AS builder
# Устанавливаем рабочую директорию
WORKDIR /app
# Устанавливаем Air для hot reload
RUN go install github.com/air-verse/air@latest
# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
# Устанавливаем зависимости
RUN go mod download
# Копируем весь исходный код проекта
COPY . .
# Этап финальной сборки (production stage)
FROM golang:1.23-alpine
# Устанавливаем рабочую директорию в контейнере
WORKDIR /app
# Копируем Air из builder stage
COPY --from=builder /go/bin/air /usr/local/bin/air
# Копируем весь исходный код и конфигурацию Air в контейнер
COPY . .
# Копируем миграции в контейнер
COPY --from=builder /app/db/migrations /migrations
# Указываем порт, который будет использовать приложение
EXPOSE 8080
# Команда для запуска Air с конфигурацией
CMD ["air", "-c", ".air.toml"]