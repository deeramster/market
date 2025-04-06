#!/bin/bash

# Запуск скрипта генерации сертификатов
echo "Генерация сертификатов..."
chmod +x generate-certs.sh
./generate-certs-source.sh
chmod +x generate-replica-certs.sh
./generate-certs-target.sh
chmod +x generate-certs-client.sh
./generate-certs-client.sh

# Запуск Docker Compose
echo "Запуск кластера Kafka..."
docker-compose up -d
