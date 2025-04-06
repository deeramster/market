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


# Генерация ключей для клиента (kafka-ui)
keytool -keystore secrets/kafka-ui/kafka.kafka-ui.keystore.jks -alias kafka-ui -validity 365 -genkey -keyalg RSA -storepass test1234 -keypass test1234 -dname "CN=kafka-ui.kafka.ssl, OU=Test, O=Company, L=City, ST=State, C=RU"
keytool -keystore secrets/kafka-ui/kafka.kafka-ui.truststore.jks -alias CARoot -import -file secrets/ca/ca-cert -storepass test1234 -keypass test1234 -noprompt

# Генерация запроса на подписание сертификата для клиента
keytool -keystore secrets/kafka-ui/kafka.kafka-ui.keystore.jks -alias kafka-ui -certreq -file secrets/kafka-ui/kafka-ui.csr -storepass test1234 -keypass test1234

# Подписание CSR клиента
openssl x509 -req -CA secrets/ca/ca-cert -CAkey secrets/ca/ca-key -in secrets/kafka-ui/kafka-ui.csr -out secrets/kafka-ui/kafka-ui-signed-cert -days 365 -CAcreateserial -passin pass:test1234

# Импорт подписанного сертификата клиента
keytool -keystore secrets/kafka-ui/kafka.kafka-ui.keystore.jks -alias CARoot -import -file secrets/ca/ca-cert -storepass test1234 -keypass test1234 -noprompt
keytool -keystore secrets/kafka-ui/kafka.kafka-ui.keystore.jks -alias kafka-ui -import -file secrets/kafka-ui/kafka-ui-signed-cert -storepass test1234 -keypass test1234 -noprompt

echo "Ктиентские сертификаты и хранилища ключей успешно созданы."
