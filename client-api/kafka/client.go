package kafka

import (
	"fmt"

	kafkago "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func CreateConsumer(config map[string]string) (*kafkago.Consumer, error) {
	consumer, err := kafkago.NewConsumer(&kafkago.ConfigMap{
		"bootstrap.servers":                     config["bootstrap.servers"],
		"group.id":                              "client-api-group",
		"auto.offset.reset":                     "earliest",
		"security.protocol":                     config["security.protocol"],
		"ssl.ca.location":                       config["ssl.ca.location"],
		"ssl.certificate.location":              config["ssl.certificate.location"],
		"ssl.key.location":                      config["ssl.key.location"],
		"ssl.key.password":                      config["ssl.key.password"],
		"ssl_endpoint_identification_algorithm": "none",
	})

	if err != nil {
		return nil, fmt.Errorf("ошибка создания consumer: %v", err)
	}

	return consumer, nil
}

func CreateProducer(config map[string]string) (*kafkago.Producer, error) {
	producer, err := kafkago.NewProducer(&kafkago.ConfigMap{
		"bootstrap.servers":                     config["bootstrap.servers"],
		"security.protocol":                     config["security.protocol"],
		"ssl.ca.location":                       config["ssl.ca.location"],
		"ssl.certificate.location":              config["ssl.certificate.location"],
		"ssl.key.location":                      config["ssl.key.location"],
		"ssl.key.password":                      config["ssl.key.password"],
		"ssl_endpoint_identification_algorithm": "none",
	})

	if err != nil {
		return nil, fmt.Errorf("ошибка создания producer: %v", err)
	}

	return producer, nil
}
