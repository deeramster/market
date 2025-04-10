package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"client-api/config"
	kafkaClient "client-api/kafka"
	"client-api/models"

	kafkago "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func SaveClientRequest(requestType, query string, config *config.Config) error {
	request := models.ClientRequest{
		Query:     query,
		Type:      requestType,
		Timestamp: time.Now(),
	}

	requestJSON, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("ошибка сериализации запроса: %v", err)
	}

	err = saveToFile(requestJSON, config.Storage)
	if err != nil {
		fmt.Printf("Предупреждение: не удалось сохранить запрос в файл: %v\n", err)
	}

	err = sendToKafka(requestJSON, config)
	if err != nil {
		fmt.Printf("Предупреждение: не удалось отправить запрос в Kafka: %v\n", err)
	}

	return nil
}

func saveToFile(data []byte, storagePath string) error {
	err := os.MkdirAll(storagePath, 0755)
	if err != nil {
		return fmt.Errorf("не удалось создать директорию хранения: %v", err)
	}

	fileName := fmt.Sprintf("request_%d.json", time.Now().UnixNano())
	filePath := filepath.Join(storagePath, fileName)

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("не удалось записать в файл: %v", err)
	}

	return nil
}

func sendToKafka(data []byte, config *config.Config) error {
	producerConfig := map[string]string{
		"bootstrap.servers":        config.KafkaSource.BootstrapServers,
		"security.protocol":        config.KafkaSource.SecurityProtocol,
		"ssl.ca.location":          config.KafkaSource.SSLCALocation,
		"ssl.certificate.location": config.KafkaSource.SSLCertificateLocation,
		"ssl.key.location":         config.KafkaSource.SSLKeyLocation,
		"ssl.key.password":         config.KafkaSource.SSLKeyPassword,
	}

	producer, err := kafkaClient.CreateProducer(producerConfig)
	if err != nil {
		return err
	}
	defer producer.Close()

	topic := config.KafkaSource.TopicClientSearch
	message := &kafkago.Message{
		TopicPartition: kafkago.TopicPartition{
			Topic:     &topic,
			Partition: kafkago.PartitionAny,
		},
		Value: data,
	}

	err = producer.Produce(message, nil)
	if err != nil {
		return fmt.Errorf("не удалось отправить сообщение в Kafka: %v", err)
	}

	producer.Flush(5000)

	return nil
}
