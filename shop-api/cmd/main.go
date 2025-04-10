package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"shop-api/internal/config"
	"shop-api/internal/generator"
	"shop-api/internal/producer"
	"shop-api/pkg/utils"
)

func main() {
	// Parse command-line flags
	var configPath string
	var generateOnly bool
	var storeID string

	flag.StringVar(&configPath, "config", "config.json", "Path to configuration file")
	flag.BoolVar(&generateOnly, "generate-only", false, "Only generate data file without sending to Kafka")
	flag.StringVar(&storeID, "store-id", "store_001", "Store ID for product generation")
	flag.Parse()

	// Load configuration
	var cfg *config.AppConfig
	var err error

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Println("Configuration file not found, using default configuration")
		cfg = config.DefaultConfig()
	} else {
		cfg, err = config.LoadConfig(configPath)
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}
	}

	// Generate products file
	log.Printf("Generating products file: %s", cfg.DataFilePath)
	productGen := generator.NewProductGenerator()
	if err := productGen.GenerateProductsFile(cfg.DataFilePath, cfg.NumProducts, storeID); err != nil {
		log.Fatalf("Failed to generate products file: %v", err)
	}
	log.Printf("Generated %d products", cfg.NumProducts)

	// If generate-only flag is set, exit
	if generateOnly {
		return
	}

	// Create Kafka producer
	log.Println("Creating Kafka producer")
	kafkaProducer, err := producer.NewKafkaProducer(&cfg.Kafka)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	// Create a channel to handle shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Read products from file and send them to Kafka
	log.Printf("Reading products from file: %s", cfg.DataFilePath)
	productChan, err := utils.ReadJSONFromFile(cfg.DataFilePath)
	if err != nil {
		log.Fatalf("Failed to read products file: %v", err)
	}

	// Process products from the channel
	log.Println("Starting to send products to Kafka")
	productsProcessed := 0

	// Create a done channel to signal completion
	done := make(chan struct{})

	go func() {
		for productData := range productChan {
			if err := kafkaProducer.SendProduct(productData); err != nil {
				log.Printf("Failed to send product: %v", err)
				continue
			}

			productsProcessed++
			if productsProcessed%10 == 0 {
				log.Printf("Processed %d products", productsProcessed)
			}
		}

		log.Printf("Finished processing all %d products", productsProcessed)
		done <- struct{}{}
	}()

	// Wait for completion or interrupt
	select {
	case <-sigChan:
		log.Println("Interrupt received, shutting down")
	case <-done:
		log.Println("All products processed")
	}

	fmt.Println("Shutdown complete")
}
