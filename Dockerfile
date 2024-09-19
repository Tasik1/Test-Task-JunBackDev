# Используем официальный образ Golang
FROM golang:latest

# Устанавливаем рабочую директорию в контейнере
WORKDIR /app

# Копируем все файлы в рабочую директорию контейнера
COPY . .

# Загружаем все зависимости
RUN go get -d -v ./...

# Устанавливаем все необходимые пакеты
RUN go install -v ./...

# Открываем порт 8080 для внешнего мира
EXPOSE 8080

# Определяем переменную окружения для строки подключения к базе данных
ENV DB_CONNECTION_STRING=postgresql://postgres:root@db/test-back-dev?sslmode=disable

# Запускаем main.go при старте контейнера
CMD ["go", "run", "main.go"]
