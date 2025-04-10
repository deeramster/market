#!/bin/bash

BASE_DIR="secrets"

DIRS=(
    "ca"
    "kafka-source-0" "kafka-source-1" "kafka-source-2"
    "kafka-target-0" "kafka-target-1" "kafka-target-2"
    "client" "kafka-ui" "schema-registry" "kafka-connect"
    "kafka-exporter" "shop-api"
)

if [ ! -d "$BASE_DIR" ]; then
    mkdir "$BASE_DIR"
    echo "Created directory: $BASE_DIR"
else
    echo "Directory $BASE_DIR already exists"
fi

for DIR in "${DIRS[@]}"; do
    FULL_PATH="$BASE_DIR/$DIR"
    if [ ! -d "$FULL_PATH" ]; then
        mkdir -p "$FULL_PATH"
        echo "Created directory: $FULL_PATH"
    else
        echo "Directory $FULL_PATH already exists"
    fi
done

echo "Generating CA (Certificate Authority)..."
openssl req -new -x509 -keyout $BASE_DIR/ca/ca-key -out $BASE_DIR/ca/ca-cert -days 365 \
    -subj '/CN=ca.kafka.ssl/OU=Test/O=Company/L=City/ST=State/C=RU' \
    -passin pass:test1234 -passout pass:test1234

echo "CA generated successfully."

generate_broker_certs() {
    local broker_type=$1
    local broker_id=$2
    local broker_name="${broker_type}-${broker_id}"

    echo "Generating certificates for $broker_name..."

    keytool -keystore $BASE_DIR/$broker_name/kafka.$broker_name.keystore.jks \
        -alias $broker_name -validity 365 -genkey -keyalg RSA \
        -storepass test1234 -keypass test1234 \
        -dname "CN=$broker_name.kafka.ssl, OU=Test, O=Company, L=City, ST=State, C=RU"

    keytool -keystore $BASE_DIR/$broker_name/kafka.$broker_name.truststore.jks \
        -alias CARoot -import -file $BASE_DIR/ca/ca-cert \
        -storepass test1234 -keypass test1234 -noprompt

    keytool -keystore $BASE_DIR/$broker_name/kafka.$broker_name.keystore.jks \
        -alias $broker_name -certreq -file $BASE_DIR/$broker_name/$broker_name.csr \
        -storepass test1234 -keypass test1234

    openssl x509 -req -CA $BASE_DIR/ca/ca-cert -CAkey $BASE_DIR/ca/ca-key \
        -in $BASE_DIR/$broker_name/$broker_name.csr \
        -out $BASE_DIR/$broker_name/$broker_name-signed-cert \
        -days 365 -CAcreateserial -passin pass:test1234

    keytool -keystore $BASE_DIR/$broker_name/kafka.$broker_name.keystore.jks \
        -alias CARoot -import -file $BASE_DIR/ca/ca-cert \
        -storepass test1234 -keypass test1234 -noprompt

    keytool -keystore $BASE_DIR/$broker_name/kafka.$broker_name.keystore.jks \
        -alias $broker_name -import -file $BASE_DIR/$broker_name/$broker_name-signed-cert \
        -storepass test1234 -keypass test1234 -noprompt

    echo "Certificates for $broker_name generated successfully."
}

generate_jks_client_certs() {
    local client_name=$1

    echo "Generating JKS certificates for $client_name..."

    keytool -keystore $BASE_DIR/$client_name/kafka.$client_name.keystore.jks \
        -alias $client_name -validity 365 -genkey -keyalg RSA \
        -storepass test1234 -keypass test1234 \
        -dname "CN=$client_name.kafka.ssl, OU=Test, O=Company, L=City, ST=State, C=RU"

    keytool -keystore $BASE_DIR/$client_name/kafka.$client_name.truststore.jks \
        -alias CARoot -import -file $BASE_DIR/ca/ca-cert \
        -storepass test1234 -keypass test1234 -noprompt

    keytool -keystore $BASE_DIR/$client_name/kafka.$client_name.keystore.jks \
        -alias $client_name -certreq -file $BASE_DIR/$client_name/$client_name.csr \
        -storepass test1234 -keypass test1234

    openssl x509 -req -CA $BASE_DIR/ca/ca-cert -CAkey $BASE_DIR/ca/ca-key \
        -in $BASE_DIR/$client_name/$client_name.csr \
        -out $BASE_DIR/$client_name/$client_name-signed-cert \
        -days 365 -CAcreateserial -passin pass:test1234

    keytool -keystore $BASE_DIR/$client_name/kafka.$client_name.keystore.jks \
        -alias CARoot -import -file $BASE_DIR/ca/ca-cert \
        -storepass test1234 -keypass test1234 -noprompt

    keytool -keystore $BASE_DIR/$client_name/kafka.$client_name.keystore.jks \
        -alias $client_name -import -file $BASE_DIR/$client_name/$client_name-signed-cert \
        -storepass test1234 -keypass test1234 -noprompt

    echo "JKS certificates for $client_name generated successfully."
}

generate_pem_client_certs() {
    local client_name=$1

    echo "Generating PEM certificates for $client_name..."

    openssl req -new -newkey rsa:4096 -nodes \
        -keyout $BASE_DIR/$client_name/$client_name-key.pem \
        -out $BASE_DIR/$client_name/$client_name.csr \
        -subj "/CN=$client_name.kafka.ssl/OU=Test/O=Company/L=City/ST=State/C=RU"

    openssl x509 -req -days 365 \
        -in $BASE_DIR/$client_name/$client_name.csr \
        -CA $BASE_DIR/ca/ca-cert \
        -CAkey $BASE_DIR/ca/ca-key \
        -CAcreateserial \
        -out $BASE_DIR/$client_name/$client_name-cert.pem \
        -passin pass:test1234

    cp $BASE_DIR/ca/ca-cert $BASE_DIR/$client_name/truststore.pem

    chmod 644 $BASE_DIR/$client_name/*.pem

    echo "PEM certificates for $client_name generated successfully."
}

for i in {0..2}; do
    generate_broker_certs "kafka-source" $i
    generate_broker_certs "kafka-target" $i
done

generate_jks_client_certs "client"
generate_jks_client_certs "kafka-ui"
generate_jks_client_certs "schema-registry"
generate_jks_client_certs "kafka-connect"

generate_pem_client_certs "kafka-exporter"
generate_pem_client_certs "shop-api"

echo "All certificates and keystores successfully created."
