package main

import (
	"flag"
	"fmt"
	"os"

	"client-api/config"
	"client-api/products"
)

func main() {
	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
	recommendCmd := flag.NewFlagSet("recommend", flag.ExitOnError)

	skuFlag := searchCmd.String("sku", "", "SKU товара для поиска")

	if len(os.Args) < 2 {
		fmt.Println("Ожидается подкоманда")
		fmt.Println("Доступные команды: search, recommend")
		os.Exit(1)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Ошибка загрузки конфигурации: %v\n", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "search":
		searchCmd.Parse(os.Args[2:])
		if *skuFlag == "" {
			fmt.Println("Укажите SKU для поиска с помощью флага --sku")
			searchCmd.PrintDefaults()
			os.Exit(1)
		}

		fmt.Printf("Поиск товара с SKU: %s\n", *skuFlag)
		product, err := products.SearchProductBySKU(*skuFlag, cfg)
		if err != nil {
			fmt.Printf("Ошибка поиска: %v\n", err)
			os.Exit(1)
		}
		products.PrintProductInfo(product)

	case "recommend":
		recommendCmd.Parse(os.Args[2:])

		fmt.Println("Получение персонализированных рекомендаций...")
		recommendations, err := products.GetPersonalizedRecommendations(cfg)
		if err != nil {
			fmt.Printf("Ошибка получения рекомендаций: %v\n", err)
			os.Exit(1)
		}
		products.PrintRecommendations(recommendations)

	default:
		fmt.Printf("Неизвестная команда: %s\n", os.Args[1])
		fmt.Println("Доступные команды: search, recommend")
		os.Exit(1)
	}
}
