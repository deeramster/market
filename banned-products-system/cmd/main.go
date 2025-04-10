package main

import (
	"fmt"
	"log"
	"os"

	"banned-products-system/internal/cli"
	"banned-products-system/internal/config"
	"banned-products-system/internal/processor"
	"banned-products-system/internal/storage"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Инициализация хранилища
	store, err := storage.New(cfg.Storage.FilePath)
	if err != nil {
		log.Fatalf("Ошибка инициализации хранилища: %v", err)
	}

	// Парсинг команд
	if len(os.Args) < 2 {
		fmt.Println("Требуется указать подкоманду: add, remove, update, list, stream")
		os.Exit(1)
	}

	// Обработка CLI-команд
	cliHandler := cli.NewHandler(store)

	// Обработка команды stream
	if proc, err := processor.New(cfg.Kafka, store); err != nil {
		log.Fatalf("Ошибка создания обработчика потоков: %v", err)
	} else {
		switch os.Args[1] {
		case "add", "remove", "update", "list":
			// Обрабатываем остальные команды
			if err := cliHandler.HandleCommand(os.Args[1:]); err != nil {
				log.Fatalf("Ошибка выполнения команды: %v", err)
			}
		case "stream":
			// Запускаем обработчик потоков
			if err := proc.Run(); err != nil {
				log.Fatalf("Ошибка запуска обработчика потоков: %v", err)
			}
		default:
			// Неизвестная команда
			fmt.Println("Неизвестная подкоманда. Используйте: add, remove, update, list, stream")
			os.Exit(1)
		}
	}
}
