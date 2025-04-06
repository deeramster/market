#!/bin/bash

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

# Генерация ключей для клиента (продюсер и консьюмер)
keytool -keystore secrets/client/kafka.client.keystore.jks -alias client -validity 365 -genkey -keyalg RSA -storepass test1234 -keypass test1234 -dname "CN=client.kafka.ssl, OU=Test, O=Company, L=City, ST=State, C=RU"
keytool -keystore secrets/client/kafka.client.truststore.jks -alias CARoot -import -file secrets/ca/ca-cert -storepass test1234 -keypass test1234 -noprompt

# Генерация запроса на подписание сертификата для клиента
keytool -keystore secrets/client/kafka.client.keystore.jks -alias client -certreq -file secrets/client/client.csr -storepass test1234 -keypass test1234

# Подписание CSR клиента
openssl x509 -req -CA secrets/ca/ca-cert -CAkey secrets/ca/ca-key -in secrets/client/client.csr -out secrets/client/client-signed-cert -days 365 -CAcreateserial -passin pass:test1234

# Импорт подписанного сертификата клиента
keytool -keystore secrets/client/kafka.client.keystore.jks -alias CARoot -import -file secrets/ca/ca-cert -storepass test1234 -keypass test1234 -noprompt
keytool -keystore secrets/client/kafka.client.keystore.jks -alias client -import -file secrets/client/client-signed-cert -storepass test1234 -keypass test1234 -noprompt

echo "Сертификаты и хранилища ключей для source кластера успешно созданы."
