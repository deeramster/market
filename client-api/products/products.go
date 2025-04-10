package products

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"client-api/config"
	kafkaClient "client-api/kafka"
	"client-api/models"
	"client-api/storage"

	kafkago "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func SearchProductBySKU(sku string, config *config.Config) (*models.Product, error) {
	consumerConfig := map[string]string{
		"bootstrap.servers":                     config.KafkaSource.BootstrapServers,
		"security.protocol":                     config.KafkaSource.SecurityProtocol,
		"ssl.ca.location":                       config.KafkaSource.SSLCALocation,
		"ssl.certificate.location":              config.KafkaSource.SSLCertificateLocation,
		"ssl.key.location":                      config.KafkaSource.SSLKeyLocation,
		"ssl.key.password":                      config.KafkaSource.SSLKeyPassword,
		"ssl_endpoint_identification_algorithm": "none",
	}

	consumer, err := kafkaClient.CreateConsumer(consumerConfig)
	if err != nil {
		return nil, err
	}
	defer consumer.Close()

	err = consumer.Subscribe(config.KafkaSource.TopicRead, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка подписки на топик: %v", err)
	}

	timeout := 10 * time.Second
	foundProduct := new(models.Product)
	found := false

	for !found {
		ev := consumer.Poll(int(timeout.Milliseconds()))
		if ev == nil {
			continue
		}

		switch e := ev.(type) {
		case *kafkago.Message:
			var product models.Product
			err := json.Unmarshal(e.Value, &product)
			if err != nil {
				fmt.Printf("Ошибка разбора данных о товаре: %v\n", err)
				continue
			}

			if product.SKU == sku {
				*foundProduct = product
				found = true
				break
			}
		case kafkago.Error:
			fmt.Printf("Ошибка Kafka: %v\n", e)
			if e.Code() == kafkago.ErrAllBrokersDown {
				return nil, fmt.Errorf("все брокеры недоступны")
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("товар с SKU %s не найден", sku)
	}

	storage.SaveClientRequest("sku", sku, config)

	return foundProduct, nil
}

func GetPersonalizedRecommendations(config *config.Config) ([]models.Product, error) {
	topHours, err := getTopHoursFromAnalytics(config)
	if err != nil {
		return nil, err
	}

	products, err := getProductsCreatedInHours(topHours, config)
	if err != nil {
		return nil, err
	}

	storage.SaveClientRequest("recommendations", "", config)

	return products, nil
}

func getTopHoursFromAnalytics(config *config.Config) ([]int, error) {
	consumerConfig := map[string]string{
		"bootstrap.servers":                     config.KafkaTarget.BootstrapServers,
		"security.protocol":                     config.KafkaTarget.SecurityProtocol,
		"ssl.ca.location":                       config.KafkaTarget.SSLCALocation,
		"ssl.certificate.location":              config.KafkaTarget.SSLCertificateLocation,
		"ssl.key.location":                      config.KafkaTarget.SSLKeyLocation,
		"ssl.key.password":                      config.KafkaTarget.SSLKeyPassword,
		"ssl_endpoint_identification_algorithm": "none",
	}

	consumer, err := kafkaClient.CreateConsumer(consumerConfig)
	if err != nil {
		return nil, err
	}
	defer consumer.Close()

	err = consumer.Subscribe(config.KafkaTarget.TopicRead, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка подписки на топик аналитики: %v", err)
	}

	timeout := 10 * time.Second
	var recommendations []models.AnalyticsRecommendation

	endTime := time.Now().Add(timeout)
	for time.Now().Before(endTime) {
		ev := consumer.Poll(100)
		if ev == nil {
			continue
		}

		switch e := ev.(type) {
		case *kafkago.Message:
			var rec models.AnalyticsRecommendation
			err := json.Unmarshal(e.Value, &rec)
			if err != nil {
				fmt.Printf("Ошибка разбора данных аналитики: %v\n", err)
				continue
			}
			recommendations = append(recommendations, rec)
		case kafkago.Error:
			fmt.Printf("Ошибка Kafka: %v\n", e)
		}
	}

	topHours := make([]int, 0, len(recommendations))
	for _, rec := range recommendations {
		topHours = append(topHours, rec.Hour)
	}

	return topHours, nil
}

func getProductsCreatedInHours(hours []int, config *config.Config) ([]models.Product, error) {
	consumerConfig := map[string]string{
		"bootstrap.servers":                     config.KafkaSource.BootstrapServers,
		"security.protocol":                     config.KafkaSource.SecurityProtocol,
		"ssl.ca.location":                       config.KafkaSource.SSLCALocation,
		"ssl.certificate.location":              config.KafkaSource.SSLCertificateLocation,
		"ssl.key.location":                      config.KafkaSource.SSLKeyLocation,
		"ssl.key.password":                      config.KafkaSource.SSLKeyPassword,
		"ssl_endpoint_identification_algorithm": "none",
	}

	consumer, err := kafkaClient.CreateConsumer(consumerConfig)
	if err != nil {
		return nil, err
	}
	defer consumer.Close()

	err = consumer.Subscribe(config.KafkaSource.TopicRead, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка подписки на топик продуктов: %v", err)
	}

	timeout := 10 * time.Second
	var matchingProducts []models.Product

	endTime := time.Now().Add(timeout)
	for time.Now().Before(endTime) {
		ev := consumer.Poll(100)
		if ev == nil {
			continue
		}

		switch e := ev.(type) {
		case *kafkago.Message:
			var product models.Product
			err := json.Unmarshal(e.Value, &product)
			if err != nil {
				fmt.Printf("Ошибка разбора данных о товаре: %v\n", err)
				continue
			}

			// Проверяем, что время создания продукта совпадает с одним из топ часов
			createdTime, err := time.Parse(time.RFC3339, product.CreatedAt)
			if err != nil {
				fmt.Printf("Ошибка разбора времени created_at: %v\n", err)
				continue
			}

			createdHour := createdTime.Hour()
			for _, hour := range hours {
				if createdHour == hour {
					matchingProducts = append(matchingProducts, product)
					break
				}
			}
		case kafkago.Error:
			fmt.Printf("Ошибка Kafka: %v\n", e)
		}
	}

	return matchingProducts, nil
}

func PrintProductInfo(product *models.Product) {
	fmt.Println("====== Информация о товаре ======")
	fmt.Printf("ID: %s\n", product.ProductID)
	fmt.Printf("Название: %s\n", product.Name)
	fmt.Printf("Описание: %s\n", product.Description)
	fmt.Printf("SKU: %s\n", product.SKU)
	fmt.Printf("Цена: %.2f %s\n", product.Price.Amount, product.Price.Currency)
	fmt.Printf("Категория: %s\n", product.Category)
	fmt.Printf("Бренд: %s\n", product.Brand)
	fmt.Printf("В наличии: %d (зарезервировано: %d)\n", product.Stock.Available, product.Stock.Reserved)
	fmt.Printf("Теги: %s\n", strings.Join(product.Tags, ", "))

	fmt.Println("Спецификации:")
	for key, value := range product.Specifications {
		fmt.Printf("  %s: %s\n", key, value)
	}

	fmt.Printf("Создан: %s\n", product.CreatedAt)
	fmt.Printf("Обновлен: %s\n", product.UpdatedAt)
	fmt.Println("================================")
}

func PrintRecommendations(products []models.Product) {
	fmt.Println("====== Персонализированные рекомендации ======")
	for i, product := range products {
		fmt.Printf("%d. %s (%.2f %s) - %s\n",
			i+1, product.Name, product.Price.Amount, product.Price.Currency, product.SKU)
	}
	fmt.Println("==============================================")
}
