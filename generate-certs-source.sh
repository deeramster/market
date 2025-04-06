#!/bin/bash

BASE_DIR="secrets"

DIRS=("kafka-source-0" "kafka-source-1" "kafka-source-2" "kafka-target-0" "kafka-target-1" "kafka-target-2" "ca" "client" "kafka-ui" "creds")

if [ ! -d "$BASE_DIR" ]; then
    mkdir "$BASE_DIR"
    echo "Создана директория: $BASE_DIR"
else
    echo "Директория $BASE_DIR уже существует"
fi

for DIR in "${DIRS[@]}"; do
    FULL_PATH="$BASE_DIR/$DIR"
    if [ ! -d "$FULL_PATH" ]; then
        mkdir "$FULL_PATH"
        echo "Создана директория: $FULL_PATH"
    else
        echo "Директория $FULL_PATH уже существует"
    fi
done

# Генерация CA (Certificate Authority)
openssl req -new -x509 -keyout secrets/ca/ca-key -out secrets/ca/ca-cert -days 365 -subj '/CN=ca.kafka.ssl/OU=Test/O=Company/L=City/ST=State/C=RU' -passin pass:test1234 -passout pass:test1234

# Генерация ключей и сертификатов для брокеров

# Broker 1
keytool -keystore secrets/kafka-source-0/kafka.kafka-source-0.keystore.jks -alias kafka-source-0 -validity 365 -genkey -keyalg RSA -storepass test1234 -keypass test1234 -dname "CN=kafka-source-0.kafka.ssl, OU=Test, O=Company, L=City, ST=State, C=RU"
keytool -keystore secrets/kafka-source-0/kafka.kafka-source-0.truststore.jks -alias CARoot -import -file secrets/ca/ca-cert -storepass test1234 -keypass test1234 -noprompt

# Broker 2
keytool -keystore secrets/kafka-source-1/kafka.kafka-source-1.keystore.jks -alias kafka-source-1 -validity 365 -genkey -keyalg RSA -storepass test1234 -keypass test1234 -dname "CN=kafka-source-1.kafka.ssl, OU=Test, O=Company, L=City, ST=State, C=RU"
keytool -keystore secrets/kafka-source-1/kafka.kafka-source-1.truststore.jks -alias CARoot -import -file secrets/ca/ca-cert -storepass test1234 -keypass test1234 -noprompt

# Broker 3
keytool -keystore secrets/kafka-source-2/kafka.kafka-source-2.keystore.jks -alias kafka-source-2 -validity 365 -genkey -keyalg RSA -storepass test1234 -keypass test1234 -dname "CN=kafka-source-2.kafka.ssl, OU=Test, O=Company, L=City, ST=State, C=RU"
keytool -keystore secrets/kafka-source-2/kafka.kafka-source-2.truststore.jks -alias CARoot -import -file secrets/ca/ca-cert -storepass test1234 -keypass test1234 -noprompt

# Генерация запросов на подписание сертификатов (CSR)
keytool -keystore secrets/kafka-source-0/kafka.kafka-source-0.keystore.jks -alias kafka-source-0 -certreq -file secrets/kafka-source-0/kafka-source-0.csr -storepass test1234 -keypass test1234
keytool -keystore secrets/kafka-source-1/kafka.kafka-source-1.keystore.jks -alias kafka-source-1 -certreq -file secrets/kafka-source-1/kafka-source-1.csr -storepass test1234 -keypass test1234
keytool -keystore secrets/kafka-source-2/kafka.kafka-source-2.keystore.jks -alias kafka-source-2 -certreq -file secrets/kafka-source-2/kafka-source-2.csr -storepass test1234 -keypass test1234

# Подписание CSR с помощью CA
openssl x509 -req -CA secrets/ca/ca-cert -CAkey secrets/ca/ca-key -in secrets/kafka-source-0/kafka-source-0.csr -out secrets/kafka-source-0/kafka-source-0-signed-cert -days 365 -CAcreateserial -passin pass:test1234
openssl x509 -req -CA secrets/ca/ca-cert -CAkey secrets/ca/ca-key -in secrets/kafka-source-1/kafka-source-1.csr -out secrets/kafka-source-1/kafka-source-1-signed-cert -days 365 -CAcreateserial -passin pass:test1234
openssl x509 -req -CA secrets/ca/ca-cert -CAkey secrets/ca/ca-key -in secrets/kafka-source-2/kafka-source-2.csr -out secrets/kafka-source-2/kafka-source-2-signed-cert -days 365 -CAcreateserial -passin pass:test1234

# Импорт подписанных сертификатов в keystore
keytool -keystore secrets/kafka-source-0/kafka.kafka-source-0.keystore.jks -alias CARoot -import -file secrets/ca/ca-cert -storepass test1234 -keypass test1234 -noprompt
keytool -keystore secrets/kafka-source-0/kafka.kafka-source-0.keystore.jks -alias kafka-source-0 -import -file secrets/kafka-source-0/kafka-source-0-signed-cert -storepass test1234 -keypass test1234 -noprompt

keytool -keystore secrets/kafka-source-1/kafka.kafka-source-1.keystore.jks -alias CARoot -import -file secrets/ca/ca-cert -storepass test1234 -keypass test1234 -noprompt
keytool -keystore secrets/kafka-source-1/kafka.kafka-source-1.keystore.jks -alias kafka-source-1 -import -file secrets/kafka-source-1/kafka-source-1-signed-cert -storepass test1234 -keypass test1234 -noprompt

keytool -keystore secrets/kafka-source-2/kafka.kafka-source-2.keystore.jks -alias CARoot -import -file secrets/ca/ca-cert -storepass test1234 -keypass test1234 -noprompt
keytool -keystore secrets/kafka-source-2/kafka.kafka-source-2.keystore.jks -alias kafka-source-2 -import -file secrets/kafka-source-2/kafka-source-2-signed-cert -storepass test1234 -keypass test1234 -noprompt
