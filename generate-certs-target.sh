#!/bin/bash

# Генерация ключей и сертификатов для брокеров
# Kafka Target 0
keytool -keystore secrets/kafka-target-0/kafka.kafka-target-0.keystore.jks -alias kafka-target-0 -validity 365 -genkey -keyalg RSA -storepass test1234 -keypass test1234 -dname "CN=kafka-target-0.kafka.ssl, OU=Test, O=Company, L=City, ST=State, C=RU"
keytool -keystore secrets/kafka-target-0/kafka.kafka-target-0.truststore.jks -alias CARoot -import -file secrets/ca/ca-cert -storepass test1234 -keypass test1234 -noprompt

# Kafka Target 1
keytool -keystore secrets/kafka-target-1/kafka.kafka-target-1.keystore.jks -alias kafka-target-1 -validity 365 -genkey -keyalg RSA -storepass test1234 -keypass test1234 -dname "CN=kafka-target-1.kafka.ssl, OU=Test, O=Company, L=City, ST=State, C=RU"
keytool -keystore secrets/kafka-target-1/kafka.kafka-target-1.truststore.jks -alias CARoot -import -file secrets/ca/ca-cert -storepass test1234 -keypass test1234 -noprompt

# Kafka Target 2
keytool -keystore secrets/kafka-target-2/kafka.kafka-target-2.keystore.jks -alias kafka-target-2 -validity 365 -genkey -keyalg RSA -storepass test1234 -keypass test1234 -dname "CN=kafka-target-2.kafka.ssl, OU=Test, O=Company, L=City, ST=State, C=RU"
keytool -keystore secrets/kafka-target-2/kafka.kafka-target-2.truststore.jks -alias CARoot -import -file secrets/ca/ca-cert -storepass test1234 -keypass test1234 -noprompt

# Генерация запросов на подписание сертификатов (CSR)
keytool -keystore secrets/kafka-target-0/kafka.kafka-target-0.keystore.jks -alias kafka-target-0 -certreq -file secrets/kafka-target-0/kafka-target-0.csr -storepass test1234 -keypass test1234
keytool -keystore secrets/kafka-target-1/kafka.kafka-target-1.keystore.jks -alias kafka-target-1 -certreq -file secrets/kafka-target-1/kafka-target-1.csr -storepass test1234 -keypass test1234
keytool -keystore secrets/kafka-target-2/kafka.kafka-target-2.keystore.jks -alias kafka-target-2 -certreq -file secrets/kafka-target-2/kafka-target-2.csr -storepass test1234 -keypass test1234

# Подписание CSR с помощью CA
openssl x509 -req -CA secrets/ca/ca-cert -CAkey secrets/ca/ca-key -in secrets/kafka-target-0/kafka-target-0.csr -out secrets/kafka-target-0/kafka-target-0-signed-cert -days 365 -CAcreateserial -passin pass:test1234
openssl x509 -req -CA secrets/ca/ca-cert -CAkey secrets/ca/ca-key -in secrets/kafka-target-1/kafka-target-1.csr -out secrets/kafka-target-1/kafka-target-1-signed-cert -days 365 -CAcreateserial -passin pass:test1234
openssl x509 -req -CA secrets/ca/ca-cert -CAkey secrets/ca/ca-key -in secrets/kafka-target-2/kafka-target-2.csr -out secrets/kafka-target-2/kafka-target-2-signed-cert -days 365 -CAcreateserial -passin pass:test1234

# Импорт подписанных сертификатов в keystore
keytool -keystore secrets/kafka-target-0/kafka.kafka-target-0.keystore.jks -alias CARoot -import -file secrets/ca/ca-cert -storepass test1234 -keypass test1234 -noprompt
keytool -keystore secrets/kafka-target-0/kafka.kafka-target-0.keystore.jks -alias kafka-target-0 -import -file secrets/kafka-target-0/kafka-target-0-signed-cert -storepass test1234 -keypass test1234 -noprompt

keytool -keystore secrets/kafka-target-1/kafka.kafka-target-1.keystore.jks -alias CARoot -import -file secrets/ca/ca-cert -storepass test1234 -keypass test1234 -noprompt
keytool -keystore secrets/kafka-target-1/kafka.kafka-target-1.keystore.jks -alias kafka-target-1 -import -file secrets/kafka-target-1/kafka-target-1-signed-cert -storepass test1234 -keypass test1234 -noprompt

keytool -keystore secrets/kafka-target-2/kafka.kafka-target-2.keystore.jks -alias CARoot -import -file secrets/ca/ca-cert -storepass test1234 -keypass test1234 -noprompt
keytool -keystore secrets/kafka-target-2/kafka.kafka-target-2.keystore.jks -alias kafka-target-2 -import -file secrets/kafka-target-2/kafka-target-2-signed-cert -storepass test1234 -keypass test1234 -noprompt

echo "Сертификаты и хранилища ключей для source кластера успешно созданы."
