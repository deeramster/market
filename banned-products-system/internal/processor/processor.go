package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"banned-products-system/internal/codec"
	"banned-products-system/internal/config"
	"banned-products-system/internal/models"
	"banned-products-system/internal/storage"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/lovoo/goka"
)

// Processor обрабатывает потоки данных из Kafka
type Processor struct {
	cfg      config.KafkaConfig
	store    *storage.Storage
	producer *kafka.Producer
	consumer *kafka.Consumer
}

// New создает новый экземпляр обработчика потоков
func New(cfg config.KafkaConfig, store *storage.Storage) (*Processor, error) {
	producer, err := NewKafkaProducer(cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания Kafka Producer: %w", err)
	}

	consumer, err := NewKafkaConsumer(cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания Kafka Consumer: %w", err)
	}

	return &Processor{
		cfg:      cfg,
		store:    store,
		producer: producer,
		consumer: consumer,
	}, nil
}

// Run запускает обработку потоков
func (p *Processor) Run() error {
	// Настраиваем логгер для Goka
	logger := log.New(os.Stderr, "[Goka] ", log.LstdFlags)

	// Создаем процессор для обработки сообщений
	g := goka.DefineGroup(goka.Group(p.cfg.GroupID),
		goka.Input(goka.Stream(p.cfg.ProductsTopic), new(codec.ProductCodec), p.processProduct),
		goka.Output(goka.Stream(p.cfg.FilteredTopic), new(codec.ProductCodec)),
	)

	// Создаем процессор с указанным логгером
	processor, err := goka.NewProcessor(strings.Split(p.cfg.Brokers, ","), g,
		goka.WithLogger(logger))
	if err != nil {
		return fmt.Errorf("ошибка создания процессора: %w", err)
	}

	// Обработка сигналов для корректного завершения
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan bool)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals
		cancel()
		close(done)
	}()

	// Запускаем обработку
	log.Println("Запуск потоковой обработки данных...")
	go func() {
		if err := processor.Run(ctx); err != nil {
			log.Fatalf("Ошибка при запуске процессора: %v", err)
		}
	}()

	fmt.Println("Поток обработки запущен. Нажмите Ctrl+C для завершения.")
	<-done
	log.Println("Завершение работы потока обработки.")

	return nil
}

// processProduct обрабатывает входящие сообщения о товарах
func (p *Processor) processProduct(ctx goka.Context, msg interface{}) {
	var product models.Product

	// Преобразуем входящее сообщение в структуру Product
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Ошибка сериализации сообщения: %v", err)
		return
	}

	if err := json.Unmarshal(data, &product); err != nil {
		log.Printf("Ошибка десериализации сообщения в структуру Product: %v", err)
		return
	}

	// Проверяем, если товар заблокирован
	if p.store.IsBanned(product.SKU) {
		log.Printf("Товар с SKU %s запрещён и будет отфильтрован", product.SKU)
		return
	}

	// Если товар не заблокирован, передаем его в отфильтрованный топик
	ctx.Emit(goka.Stream(p.cfg.FilteredTopic), ctx.Key(), msg)
	log.Printf("Товар с SKU %s прошёл фильтрацию и отправлен далее", product.SKU)
}
