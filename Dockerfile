# syntax=docker/dockerfile:1.4

# 1-й этап: Сборка
FROM golang:1.21 AS builder

WORKDIR /app

# Скопировать go.mod, go.sum и скачать зависимости
COPY go.mod go.sum ./
RUN go mod download

# Скопировать остальной код
COPY . .

# Собрать бинарник
RUN go build -o mybot main.go

# 2-й этап: Runtime (запуск)
FROM alpine:latest

# Папка, в которой будет лежать бинарник
WORKDIR /app

# Копируем бинарник из 1-го этапа
COPY --from=builder /app/mybot /app/mybot

# Опционально устанавливаем сертификаты, если нужен HTTPS или webhook
RUN apk add --no-cache ca-certificates

# При запуске контейнера выполняется эта команда
CMD ["/app/mybot"]
