#!/bin/bash

# Применим миграции
# ВАЖНО! Для локального запуска вызвать из корня проекта, предварительно подгрузив переменные окружения в терминал командой  `set -a; source <путь к .env файлу>; set +a`
# Установка goose: go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir ./migrations postgres "user=$COMMUNITY_SERVICE_POSTGRES_USER password=$COMMUNITY_SERVICE_POSTGRES_PASSWORD dbname=$COMMUNITY_SERVICE_POSTGRES_DB host=$COMMUNITY_SERVICE_POSTGRES_HOST port=$COMMUNITY_SERVICE_POSTGRES_PORT sslmode=disable" up 1