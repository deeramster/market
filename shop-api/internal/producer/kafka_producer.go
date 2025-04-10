package producer

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"

	"shop-api/internal/config"
	"shop-api/internal/models"
)

// KafkaProducer handles Kafka message production
type KafkaProducer struct {
	producer   *kafka.Producer
	serializer *jsonschema.Serializer
	topic      string
}

// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer(cfg *config.KafkaConfig) (*KafkaProducer, error) {
	// Configure the Kafka producer
	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers":                     cfg.BootstrapServers,
		"security.protocol":                     cfg.SecurityProtocol,
		"ssl.ca.location":                       cfg.SslCaLocation,
		"ssl.certificate.location":              cfg.SslCertificateLocation,
		"ssl.key.location":                      cfg.SslKeyLocation,
		"ssl.key.password":                      cfg.SslKeyPassword,
		"ssl.endpoint.identification.algorithm": "none",
	}

	producer, err := kafka.NewProducer(kafkaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	// Configure Schema Registry client
	srClient, err := schemaregistry.NewClient(schemaregistry.NewConfig(cfg.SchemaRegistryURL))
	if err != nil {
		producer.Close()
		return nil, fmt.Errorf("failed to create Schema Registry client: %w", err)
	}

	// Create JSON serializer
	serializer, err := jsonschema.NewSerializer(srClient, serde.ValueSerde, jsonschema.NewSerializerConfig())
	if err != nil {
		producer.Close()
		return nil, fmt.Errorf("failed to create serializer: %w", err)
	}

	// Start monitoring delivery reports in a goroutine
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Failed to deliver message: %v\n", ev.TopicPartition.Error)
				} else {
					log.Printf("Successfully produced message to topic %s partition [%d] @ offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			}
		}
	}()

	return &KafkaProducer{
		producer:   producer,
		serializer: serializer,
		topic:      cfg.Topic,
	}, nil
}

// SendProduct sends a product to Kafka
func (kp *KafkaProducer) SendProduct(productData []byte) error {
	// Parse JSON into a Product struct
	var product models.Product
	if err := json.Unmarshal(productData, &product); err != nil {
		return fmt.Errorf("failed to unmarshal product: %w", err)
	}

	// Serialize the product using the schema registry
	payload, err := kp.serializer.Serialize(kp.topic, &product)
	if err != nil {
		return fmt.Errorf("failed to serialize product: %w", err)
	}

	// Send the message to Kafka
	err = kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kp.topic, Partition: kafka.PartitionAny},
		Key:            []byte(product.ProductID),
		Value:          payload,
		Headers:        []kafka.Header{{Key: "source", Value: []byte("shop-api")}},
	}, nil)

	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	return nil
}

// Close closes the Kafka producer
func (kp *KafkaProducer) Close() {
	// Flush any remaining messages
	kp.producer.Flush(15 * 1000)
	kp.producer.Close()
}
