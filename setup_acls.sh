#!/bin/bash

set -euo pipefail

echo "🔐 Установка ACL для SOURCE кластера..."

docker-compose exec kafka-source-0 kafka-acls.sh --bootstrap-server kafka-source-0:9092 \
  --add --allow-principal "User:C=RU,ST=State,L=City,O=Company,OU=Test,CN=kafka-connect.kafka.ssl" \
  --operation All --topic '*' --group '*' --cluster \
  --command-config /bitnami/kafka/config/client-ssl.properties


docker-compose exec kafka-source-0 kafka-acls.sh --bootstrap-server kafka-source-0:9092 \
  --add --allow-principal "User:C=RU,ST=State,L=City,O=Company,OU=Test,CN=shop-api.kafka.ssl" \
  --operation Create --operation Alter --operation Write --operation Read --operation Describe \
  --topic products \
  --command-config /bitnami/kafka/config/client-ssl.properties

docker-compose exec kafka-source-0 kafka-acls.sh --bootstrap-server kafka-source-0:9092 \
  --add --allow-principal "User:C=RU,ST=State,L=City,O=Company,OU=Test,CN=client-api.kafka.ssl" \
  --operation Read --topic products \
  --command-config /bitnami/kafka/config/client-ssl.properties

docker-compose exec kafka-source-0 kafka-acls.sh --bootstrap-server kafka-source-0:9092 \
  --add --allow-principal "User:C=RU,ST=State,L=City,O=Company,OU=Test,CN=shop-api.kafka.ssl" \
  --operation Create --operation Alter --operation Write --operation Read --operation Describe \
  --topic filtered-products \
  --command-config /bitnami/kafka/config/client-ssl.properties

docker-compose exec kafka-source-0 kafka-acls.sh --bootstrap-server kafka-source-0:9092 \
  --add --allow-principal "User:C=RU,ST=State,L=City,O=Company,OU=Test,CN=client-api.kafka.ssl" \
  --operation Read --topic filtered-products \
  --command-config /bitnami/kafka/config/client-ssl.properties

docker-compose exec kafka-source-0 kafka-acls.sh --bootstrap-server kafka-source-0:9092 \
  --add --allow-principal "User:C=RU,ST=State,L=City,O=Company,OU=Test,CN=client-api.kafka.ssl" \
  --operation Read --group '*' \
  --command-config /bitnami/kafka/config/client-ssl.properties

docker-compose exec kafka-source-0 kafka-acls.sh --bootstrap-server kafka-source-0:9092 \
  --add --allow-principal "User:C=RU,ST=State,L=City,O=Company,OU=Test,CN=client-api.kafka.ssl" \
  --operation Write --topic client-activity \
  --command-config /bitnami/kafka/config/client-ssl.properties

docker-compose exec kafka-source-0 kafka-acls.sh --bootstrap-server kafka-source-0:9092 \
  --add --allow-principal "User:C=RU,ST=State,L=City,O=Company,OU=Test,CN=schema-registry.kafka.ssl" \
  --operation Read --operation Write --operation Create \
  --topic __schemas \
  --command-config /bitnami/kafka/config/client-ssl.properties

docker-compose exec kafka-source-0 kafka-acls.sh --bootstrap-server kafka-source-0:9092 \
  --add --allow-principal "User:C=RU,ST=State,L=City,O=Company,OU=Test,CN=schema-registry.kafka.ssl" \
  --cluster --operation ClusterAction \
  --command-config /bitnami/kafka/config/client-ssl.properties

docker-compose exec kafka-source-0 kafka-acls.sh --bootstrap-server kafka-source-0:9092 \
  --add --allow-principal "User:C=RU,ST=State,L=City,O=Company,OU=Test,CN=kafka-ui.kafka.ssl" \
  --operation Read --operation Describe --topic '*' \
  --command-config /bitnami/kafka/config/client-ssl.properties

docker-compose exec kafka-source-0 kafka-acls.sh --bootstrap-server kafka-source-0:9092 \
  --add --allow-principal "User:C=RU,ST=State,L=City,O=Company,OU=Test,CN=kafka-ui.kafka.ssl" \
  --operation Read --operation Describe --group '*' \
  --command-config /bitnami/kafka/config/client-ssl.properties

  # Для SOURCE кластера
  docker-compose exec kafka-source-0 kafka-acls.sh --bootstrap-server kafka-source-0:9092 \
    --add --allow-principal "User:C=RU,ST=State,L=City,O=Company,OU=Test,CN=kafka-ui.kafka.ssl" \
    --operation Describe --operation Read --topic '*' \
    --command-config /bitnami/kafka/config/client-ssl.properties

  docker-compose exec kafka-source-0 kafka-acls.sh --bootstrap-server kafka-source-0:9092 \
    --add --allow-principal "User:C=RU,ST=State,L=City,O=Company,OU=Test,CN=kafka-ui.kafka.ssl" \
    --operation Describe --operation Read --group '*' \
    --command-config /bitnami/kafka/config/client-ssl.properties

  docker-compose exec kafka-source-0 kafka-acls.sh --bootstrap-server kafka-source-0:9092 \
    --add --allow-principal "User:C=RU,ST=State,L=City,O=Company,OU=Test,CN=kafka-ui.kafka.ssl" \
    --operation Describe --cluster \
    --command-config /bitnami/kafka/config/client-ssl.properties

  # Для TARGET кластера (аналогично)
  docker-compose exec kafka-target-0 kafka-acls.sh --bootstrap-server kafka-target-0:9092 \
    --add --allow-principal "User:C=RU,ST=State,L=City,O=Company,OU=Test,CN=kafka-ui.kafka.ssl" \
    --operation Describe --operation Read --topic '*' \
    --command-config /bitnami/kafka/config/client-ssl.properties

  docker-compose exec kafka-target-0 kafka-acls.sh --bootstrap-server kafka-target-0:9092 \
    --add --allow-principal "User:C=RU,ST=State,L=City,O=Company,OU=Test,CN=kafka-ui.kafka.ssl" \
    --operation Describe --operation Read --group '*' \
    --command-config /bitnami/kafka/config/client-ssl.properties

  docker-compose exec kafka-target-0 kafka-acls.sh --bootstrap-server kafka-target-0:9092 \
    --add --allow-principal "User:C=RU,ST=State,L=City,O=Company,OU=Test,CN=kafka-ui.kafka.ssl" \
    --operation Describe --cluster \
    --command-config /bitnami/kafka/config/client-ssl.properties

echo "✅ SOURCE кластер: ACL установлены"

echo "🔐 Установка ACL для TARGET кластера..."

docker-compose exec kafka-target-0 kafka-acls.sh --bootstrap-server kafka-target-0:9092 \
  --add --allow-principal "User:C=RU,ST=State,L=City,O=Company,OU=Test,CN=kafka-connect.kafka.ssl" \
  --operation All --topic '*' --group '*' --cluster \
  --command-config /bitnami/kafka/config/client-ssl.properties

  docker-compose exec kafka-target-0 kafka-acls.sh --bootstrap-server kafka-target-0:9092 \
    --add --allow-principal "User:C=RU,ST=State,L=City,O=Company,OU=Test,CN=kafka-ui.kafka.ssl" \
    --operation Read --operation Describe --topic '*' \
    --command-config /bitnami/kafka/config/client-ssl.properties

  docker-compose exec kafka-target-0 kafka-acls.sh --bootstrap-server kafka-target-0:9092 \
    --add --allow-principal "User:C=RU,ST=State,L=City,O=Company,OU=Test,CN=kafka-ui.kafka.ssl" \
    --operation Read --operation Describe --group '*' \
    --command-config /bitnami/kafka/config/client-ssl.properties

docker-compose exec kafka-target-0 kafka-acls.sh --bootstrap-server kafka-target-0:9092 \
  --add --allow-principal "User:CN=kafka-target-0.kafka.ssl,OU=Test,O=Company,L=City,ST=State,C=RU" \
  --allow-principal "User:CN=kafka-target-1.kafka.ssl,OU=Test,O=Company,L=City,ST=State,C=RU" \
  --operation All  \
  --topic '*'  \
  --group '*'  \
  --cluster \
  --command-config /bitnami/kafka/config/client-ssl.properties

  docker-compose exec kafka-target-0 kafka-acls.sh --bootstrap-server kafka-target-0:9092 \
    --add --allow-principal "User:C=RU,ST=State,L=City,O=Company,OU=Test,CN=kafka-ui.kafka.ssl" \
    --operation Describe --operation Read --topic '*' \
    --command-config /bitnami/kafka/config/client-ssl.properties

  docker-compose exec kafka-target-0 kafka-acls.sh --bootstrap-server kafka-target-0:9092 \
    --add --allow-principal "User:C=RU,ST=State,L=City,O=Company,OU=Test,CN=kafka-ui.kafka.ssl" \
    --operation Describe --operation Read --group '*' \
    --command-config /bitnami/kafka/config/client-ssl.properties

  docker-compose exec kafka-target-0 kafka-acls.sh --bootstrap-server kafka-target-0:9092 \
    --add --allow-principal "User:C=RU,ST=State,L=City,O=Company,OU=Test,CN=kafka-ui.kafka.ssl" \
    --operation Describe --cluster \
    --command-config /bitnami/kafka/config/client-ssl.properties

echo "✅ TARGET кластер: ACL установлены"

echo ""
echo "🔍 Проверка ACL для SOURCE кластера:"
docker-compose exec kafka-source-0 kafka-acls.sh --bootstrap-server kafka-source-0:9092 --list --command-config /bitnami/kafka/config/client-ssl.properties

echo ""
echo "🔍 Проверка ACL для TARGET кластера:"
docker-compose exec kafka-target-0 kafka-acls.sh --bootstrap-server kafka-target-0:9092 --list --command-config /bitnami/kafka/config/client-ssl.properties

echo ""
echo "Все ACL успешно применены и проверены!"
