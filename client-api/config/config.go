package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	KafkaSource struct {
		BootstrapServers       string `json:"bootstrap_servers"`
		SchemaRegistryURL      string `json:"schema_registry_url"`
		TopicRead              string `json:"topic_read"`
		TopicClientSearch      string `json:"topic_write"`
		SecurityProtocol       string `json:"security_protocol"`
		SSLCALocation          string `json:"ssl_ca_location"`
		SSLCertificateLocation string `json:"ssl_certificate_location"`
		SSLKeyLocation         string `json:"ssl_key_location"`
		SSLKeyPassword         string `json:"ssl_key_password"`
	} `json:"kafka-source"`
	KafkaTarget struct {
		BootstrapServers       string `json:"bootstrap_servers"`
		SchemaRegistryURL      string `json:"schema_registry_url"`
		TopicRead              string `json:"topic_read"`
		SecurityProtocol       string `json:"security_protocol"`
		SSLCALocation          string `json:"ssl_ca_location"`
		SSLCertificateLocation string `json:"ssl_certificate_location"`
		SSLKeyLocation         string `json:"ssl_key_location"`
		SSLKeyPassword         string `json:"ssl_key_password"`
	} `json:"kafka-target"`
	Storage string `json:"storage"`
}

func LoadConfig() (*Config, error) {
	configPath := filepath.Join("config", "config.json")
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла конфигурации: %v", err)
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, fmt.Errorf("ошибка разбора конфигурации: %v", err)
	}

	return &config, nil
}
