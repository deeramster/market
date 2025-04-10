#!/bin/bash
set -e # Exit immediately if a command exits with a non-zero status

PROJECT_NAME="shop-api"

if [ ! -d "$PROJECT_NAME" ]; then
    mkdir -p "$PROJECT_NAME"
else
    echo "Директория '$PROJECT_NAME' уже существует."
    exit 1
fi

cd "$PROJECT_NAME"

mkdir -p cmd internal/config internal/models internal/producer internal/generator pkg/utils

touch cmd/main.go
touch internal/config/config.go
touch internal/models/product.go
touch internal/producer/kafka_producer.go
touch internal/generator/file_generator.go
touch pkg/utils/utils.go
touch go.mod

cd ..

echo "Структура проекта '$PROJECT_NAME' успешно создана."
