package config

import (
	"encoding/json"
	"os"
	"strings"
)

// Config contains the entire application configuration
type Config struct {
	Kafka   KafkaConfig   `json:"kafka"`
	Storage StorageConfig `json:"storage"`
}

// KafkaConfig contains settings for Kafka
type KafkaConfig struct {
	Brokers                            string `json:"brokers"`
	ProductsTopic                      string `json:"products_topic"`
	FilteredTopic                      string `json:"filtered_topic"`
	GroupID                            string `json:"group_id"`
	SecurityProtocol                   string `json:"security_protocol"`
	SSLCA                              string `json:"ssl_ca"`
	SSLCert                            string `json:"ssl_cert"`
	SSLKey                             string `json:"ssl_key"`
	SSLKeyPassword                     string `json:"ssl_key_password"`
	SSLEndpointIdentificationAlgorithm string `json:"ssl_endpoint_identification_algorithm"`
}

// IsSSLEnabled returns true if SSL/TLS is enabled
func (k *KafkaConfig) IsSSLEnabled() bool {
	return strings.ToLower(k.SecurityProtocol) == "ssl"
}

// StorageConfig contains settings for data storage
type StorageConfig struct {
	FilePath string `json:"config.json"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		Kafka: KafkaConfig{
			Brokers:                            "kafka-source-0:9092",
			ProductsTopic:                      "products",
			FilteredTopic:                      "filtered-products",
			GroupID:                            "products-filter-group",
			SecurityProtocol:                   "plaintext",
			SSLCA:                              "",
			SSLCert:                            "",
			SSLKey:                             "",
			SSLKeyPassword:                     "",
			SSLEndpointIdentificationAlgorithm: "none",
		},
		Storage: StorageConfig{
			FilePath: "/app/storage/data/banned_products.json",
		},
	}
}

// Load loads configuration from file or returns default values
func Load() (Config, error) {
	cfg := DefaultConfig()

	configPath := "/app/storage/config.json"
	if _, err := os.Stat(configPath); err == nil {
		file, err := os.ReadFile(configPath)
		if err != nil {
			return cfg, err
		}

		if err := json.Unmarshal(file, &cfg); err != nil {
			return cfg, err
		}
	}

	return cfg, nil
}
