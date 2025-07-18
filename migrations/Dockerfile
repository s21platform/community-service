# Используем официальный образ Golang как базовый
FROM golang:1.24-alpine

# Установим необходимые зависимости
RUN apk add --no-cache bash gcc musl-dev postgresql-dev

# Установим Goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Создадим рабочую директорию
WORKDIR /app

# Скопируем файлы миграций
COPY . ./migrations

# Скопируем entrypoint скрипт !!!! данных файлов нет
COPY entrypoint_up_all.sh /app/entrypoint_up_all.sh
COPY entrypoint_up_one.sh /app/entrypoint_up_one.sh
COPY entrypoint_down_one.sh /app/entrypoint_down_one.sh

# Сделаем entrypoint скрипт исполняемым
RUN chmod +x /app/entrypoint_up_all.sh
RUN chmod +x /app/entrypoint_up_one.sh
RUN chmod +x /app/entrypoint_down_one.sh

# Установим entrypoint
#ENTRYPOINT ["/app/entrypoint_up_all.sh"]