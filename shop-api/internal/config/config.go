package config

import (
	"encoding/json"
	"os"
)

// KafkaConfig contains Kafka configuration
type KafkaConfig struct {
	BootstrapServers       string `json:"bootstrap_servers"`
	SchemaRegistryURL      string `json:"schema_registry_url"`
	Topic                  string `json:"topic"`
	SecurityProtocol       string `json:"security_protocol"`
	SslCaLocation          string `json:"ssl_ca_location"`
	SslCertificateLocation string `json:"ssl_certificate_location"`
	SslKeyLocation         string `json:"ssl_key_location"`
	SslKeyPassword         string `json:"ssl_key_password"`
}

// AppConfig contains application configuration
type AppConfig struct {
	DataFilePath string      `json:"data_file_path"`
	NumProducts  int         `json:"num_products"`
	Kafka        KafkaConfig `json:"kafka"`
}

// LoadConfig loads configuration from file
func LoadConfig(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// DefaultConfig returns default configuration
func DefaultConfig() *AppConfig {
	return &AppConfig{
		DataFilePath: "/app/data/products.json",
		NumProducts:  100,
		Kafka: KafkaConfig{
			BootstrapServers:       "kafka-source-0:9092,kafka-source-1:9092",
			SchemaRegistryURL:      "http://schema-registry:8081",
			Topic:                  "products",
			SecurityProtocol:       "SSL",
			SslCaLocation:          "/app/certs/truststore.pem",
			SslCertificateLocation: "/app/certs/shop-api-cert.pem",
			SslKeyLocation:         "/app/certs/shop-api-key.pem",
			SslKeyPassword:         "test1234",
		},
	}
}
