package producer

import (
	"context"
	"fmt"

	"spark-app/config"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"
)

type AnalyticsProducer struct {
	producer   *kafka.Producer
	serializer *jsonschema.Serializer
	topic      string
}

type Recommendation struct {
	Hour  int `json:"hour"`
	Count int `json:"count"`
}

func NewAnalyticsProducer(cfg kafka.ConfigMap) (*AnalyticsProducer, error) {
	schemaRegistryClient, err := schemaregistry.NewClient(
		schemaregistry.NewConfig(
			config.SchemaRegistryURL(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания клиента schema-registry: %w", err)
	}

	serializerConfig := jsonschema.NewSerializerConfig()
	serializerConfig.AutoRegisterSchemas = true

	serializer, err := jsonschema.NewSerializer(
		schemaRegistryClient,
		serde.ValueSerde,
		serializerConfig,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания сериализатора: %w", err)
	}

	p, err := kafka.NewProducer(&cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания продюсера: %w", err)
	}

	return &AnalyticsProducer{
		producer:   p,
		serializer: serializer,
		topic:      config.AnalyticsTopic(),
	}, nil
}

func (a *AnalyticsProducer) SendRecommendation(ctx context.Context, rec Recommendation) error {
	payload, err := a.serializer.Serialize(a.topic, &rec)
	if err != nil {
		return fmt.Errorf("ошибка сериализации рекомендации: %w", err)
	}

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &a.topic, Partition: kafka.PartitionAny},
		Value:          payload,
	}

	deliveryChan := make(chan kafka.Event, 1)
	err = a.producer.Produce(msg, deliveryChan)
	if err != nil {
		return fmt.Errorf("ошибка отправки сообщения: %w", err)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case e := <-deliveryChan:
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			return fmt.Errorf("ошибка при доставке: %w", m.TopicPartition.Error)
		}
	}

	return nil
}

func (a *AnalyticsProducer) Close() {
	a.producer.Close()
}
