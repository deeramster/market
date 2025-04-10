#!/bin/bash

# Запуск скрипта генерации сертификатов
echo "Генерация сертификатов..."
chmod +x generate-certs.sh
./generate-certs-source.sh

# Запуск Docker Compose
echo "Запуск кластера Kafka..."
docker-compose up -d


# Настройка ACL
echo "Настройка ACL..."
chmod +x setup_acls.sh
./setup_acls.sh
