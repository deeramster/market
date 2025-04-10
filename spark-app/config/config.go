package config

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func KafkaConfig() *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers":                     "kafka-target-0:9092,kafka-target-1:9092",
		"security.protocol":                     "SSL",
		"ssl.ca.location":                       "certs/truststore.pem",
		"ssl.certificate.location":              "certs/shop-api-cert.pem",
		"ssl.key.location":                      "certs/shop-api-key.pem",
		"ssl.key.password":                      "test1234",
		"group.id":                              "analytics-consumer-group",
		"ssl.endpoint.identification.algorithm": "none",
		"auto.offset.reset":                     "earliest",
		"fetch.min.bytes":                       1024,
	}
}

func SchemaRegistryURL() string {
	return "http://schema-registry-target:8881"
}

func SparkURL() string {
	return "sc://spark:15002"
}

func HDFSNamenode() string {
	return "hadoop-namenode:9000"
}

func KafkaTopic() string {
	return "filtered_products"
}

func AnalyticsTopic() string {
	return "analytics"
}

func HDFSDataPath() string {
	return "/data/filtered_products"
}
