package processor

import (
	"banned-products-system/internal/config"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// NewKafkaProducer создает новый Kafka Producer с SSL-настройками
func NewKafkaProducer(cfg config.KafkaConfig) (*kafka.Producer, error) {
	kafkaCfg := &kafka.ConfigMap{
		"bootstrap.servers":                     cfg.Brokers,
		"security.protocol":                     "ssl",
		"ssl.ca.location":                       cfg.SSLCA,
		"ssl.certificate.location":              cfg.SSLCert,
		"ssl.key.location":                      cfg.SSLKey,
		"ssl.key.password":                      cfg.SSLKeyPassword,
		"ssl.endpoint.identification.algorithm": cfg.SSLEndpointIdentificationAlgorithm,
	}

	// Создание Kafka Producer
	producer, err := kafka.NewProducer(kafkaCfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания Kafka Producer: %w", err)
	}

	log.Println("Kafka Producer создан успешно")
	return producer, nil
}

// NewKafkaConsumer создает новый Kafka Consumer с SSL-настройками
func NewKafkaConsumer(cfg config.KafkaConfig) (*kafka.Consumer, error) {
	kafkaCfg := &kafka.ConfigMap{
		"bootstrap.servers":                     cfg.Brokers,
		"security.protocol":                     "ssl",
		"ssl.ca.location":                       cfg.SSLCA,
		"ssl.certificate.location":              cfg.SSLCert,
		"ssl.key.location":                      cfg.SSLKey,
		"ssl.key.password":                      cfg.SSLKeyPassword,
		"ssl.endpoint.identification.algorithm": cfg.SSLEndpointIdentificationAlgorithm,
		"group.id":                              cfg.GroupID,
		"auto.offset.reset":                     "earliest",
	}

	// Создание Kafka Consumer
	consumer, err := kafka.NewConsumer(kafkaCfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания Kafka Consumer: %w", err)
	}

	// Подписка на топики
	err = consumer.Subscribe(cfg.ProductsTopic, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка подписки на топик: %w", err)
	}

	log.Println("Kafka Consumer создан и подписан на топик:", cfg.ProductsTopic)
	return consumer, nil
}
