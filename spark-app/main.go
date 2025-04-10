package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"spark-app/analytics"
	"spark-app/config"
	"spark-app/consumer"
	"spark-app/producer"
)

func main() {
	analyticsModePtr := flag.Bool("analytics", false, "Запустить только аналитику")
	consumerModePtr := flag.Bool("consumer", false, "Запустить только consumer")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("Получен сигнал завершения...")
		cancel()
	}()

	if *analyticsModePtr {
		runAnalytics(ctx)
	} else if *consumerModePtr {
		runConsumer(ctx)
	} else {
		go runAnalytics(ctx)
		runConsumer(ctx)
	}
}

func runAnalytics(ctx context.Context) {
	log.Println("Запуск аналитики...")

	recommendations, err := analytics.AnalyzeCreatedAtHours()
	if err != nil {
		log.Fatalf("Ошибка аналитики: %v", err)
	}

	analyticsProducer, err := producer.NewAnalyticsProducer(*config.KafkaConfig())
	if err != nil {
		log.Fatalf("Ошибка создания продюсера: %v", err)
	}
	defer analyticsProducer.Close()

	count := 0
	for _, rec := range recommendations {
		if count >= 3 {
			break
		}

		err := analyticsProducer.SendRecommendation(ctx, producer.Recommendation{
			Hour:  int(rec.Hour),
			Count: int(rec.Count),
		})
		if err != nil {
			log.Printf("Ошибка отправки: %v", err)
		} else {
			log.Printf("Рекомендация отправлена: час %d, количество %d", rec.Hour, rec.Count)
		}
		count++
	}

	log.Println("Аналитика завершена и результаты отправлены.")
}

func runConsumer(ctx context.Context) {
	log.Println("Запуск consumer...")

	cons, err := consumer.NewConsumer(config.KafkaConfig())
	if err != nil {
		log.Fatalf("Ошибка создания consumer: %v", err)
	}
	defer cons.Close()

	err = cons.Run(ctx)
	if err != nil {
		log.Fatalf("Ошибка в работе consumer: %v", err)
	}
}
