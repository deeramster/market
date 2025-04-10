package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"

	"spark-app/config"
	"spark-app/hdfs"
)

type Consumer struct {
	kafkaConsumer *kafka.Consumer
	schemaClient  schemaregistry.Client
	deserializer  *jsonschema.Deserializer
	hdfsWriter    *hdfs.Writer
	topic         string
}

func NewConsumer(kafkaConfig *kafka.ConfigMap) (*Consumer, error) {
	consumer, err := kafka.NewConsumer(kafkaConfig)
	if err != nil {
		return nil, err
	}

	err = consumer.Subscribe(config.KafkaTopic(), nil)
	if err != nil {
		consumer.Close()
		return nil, err
	}

	schemaClient, err := schemaregistry.NewClient(
		schemaregistry.NewConfig(
			config.SchemaRegistryURL(),
		),
	)
	if err != nil {
		consumer.Close()
		return nil, err
	}

	deserializer, err := jsonschema.NewDeserializer(
		schemaClient,
		serde.ValueSerde,
		jsonschema.NewDeserializerConfig(),
	)
	if err != nil {
		consumer.Close()
		return nil, err
	}

	hdfsWriter, err := hdfs.NewWriter()
	if err != nil {
		consumer.Close()
		return nil, err
	}

	return &Consumer{
		kafkaConsumer: consumer,
		schemaClient:  schemaClient,
		deserializer:  deserializer,
		hdfsWriter:    hdfsWriter,
		topic:         config.KafkaTopic(),
	}, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	log.Println("Consumer запущен. Чтение сообщений...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Завершение по сигналу.")
			return nil
		default:
			msg, err := c.kafkaConsumer.ReadMessage(100)
			if err != nil {
				if err.(kafka.Error).Code() == kafka.ErrTimedOut {
					continue
				}
				log.Printf("Ошибка чтения сообщения: %v", err)
				continue
			}

			var data map[string]interface{}
			topicStr := ""
			if msg.TopicPartition.Topic != nil {
				topicStr = *msg.TopicPartition.Topic
			} else {
				topicStr = c.topic
			}

			err = c.deserializer.DeserializeInto(topicStr, msg.Value, &data)
			if err != nil {
				log.Printf("Ошибка десериализации: %v", err)
				continue
			}

			jsonData, _ := json.MarshalIndent(data, "", "  ")
			log.Printf("Валидированное сообщение:\n%s\n", string(jsonData))

			err = c.hdfsWriter.WriteJSONLine(jsonData)
			if err != nil {
				log.Printf("Ошибка записи в HDFS: %v", err)
			} else {
				log.Printf("Сообщение записано в HDFS")
			}
		}
	}
}

func (c *Consumer) Close() {
	c.kafkaConsumer.Close()
	if c.hdfsWriter != nil {
		err := c.hdfsWriter.Close()
		if err != nil {
			log.Printf("Ошибка закрытия HDFS: %v", err)
		}
	}
}
