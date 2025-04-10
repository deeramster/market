package generator

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"shop-api/internal/models"
	"shop-api/pkg/utils"
)

// ProductGenerator generates product data
type ProductGenerator struct {
	brands     []string
	categories []string
	tags       [][]string
}

// NewProductGenerator creates a new product generator
func NewProductGenerator() *ProductGenerator {
	return &ProductGenerator{
		brands:     []string{"XYZ", "ABC", "TechGiant", "SmartLife", "FutureTech"},
		categories: []string{"Электроника", "Бытовая техника", "Компьютеры", "Смартфоны", "Аксессуары"},
		tags: [][]string{
			{"умные часы", "гаджеты", "технологии"},
			{"смартфон", "телефон", "мобильный"},
			{"ноутбук", "компьютер", "работа"},
			{"наушники", "аудио", "музыка"},
			{"камера", "фото", "видео"},
		},
	}
}

// GenerateProduct generates a single product
func (pg *ProductGenerator) GenerateProduct(index int, storeID string) models.Product {
	now := time.Now()
	r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(index)))

	brand := pg.brands[r.Intn(len(pg.brands))]
	category := pg.categories[r.Intn(len(pg.categories))]
	tagSet := pg.tags[r.Intn(len(pg.tags))]
	productID := fmt.Sprintf("%d", 10000+index)

	return models.Product{
		ProductID:   productID,
		Name:        fmt.Sprintf("%s %s %d", brand, category, index),
		Description: fmt.Sprintf("%s с расширенными функциональными возможностями", category),
		Price: models.Price{
			Amount:   float64(1000+r.Intn(9000)) + r.Float64(),
			Currency: "RUB",
		},
		Category: category,
		Brand:    brand,
		Stock: models.Stock{
			Available: 50 + r.Intn(200),
			Reserved:  r.Intn(50),
		},
		SKU:  fmt.Sprintf("%s-%s", brand, productID),
		Tags: tagSet,
		Images: []models.Image{
			{
				URL: fmt.Sprintf("https://example.com/images/product%d.jpg", index),
				Alt: fmt.Sprintf("%s %s - вид спереди", brand, category),
			},
			{
				URL: fmt.Sprintf("https://example.com/images/product%d_side.jpg", index),
				Alt: fmt.Sprintf("%s %s - вид сбоку", brand, category),
			},
		},
		Specifications: models.Specifications{
			Weight:          fmt.Sprintf("%dg", 50+r.Intn(450)),
			Dimensions:      fmt.Sprintf("%dmm x %dmm x %dmm", 10+r.Intn(90), 10+r.Intn(90), 5+r.Intn(25)),
			BatteryLife:     fmt.Sprintf("%d hours", 12+r.Intn(36)),
			WaterResistance: fmt.Sprintf("IP%d", 54+r.Intn(14)),
		},
		CreatedAt: now.Add(-time.Duration(r.Intn(30)) * 24 * time.Hour),
		UpdatedAt: now.Add(-time.Duration(r.Intn(10)) * 24 * time.Hour),
		Index:     "products",
		StoreID:   storeID,
	}
}

// GenerateProductsFile generates a file with product data
func (pg *ProductGenerator) GenerateProductsFile(filename string, count int, storeID string) error {
	products := make([]models.Product, count)

	for i := 0; i < count; i++ {
		products[i] = pg.GenerateProduct(i, storeID)
	}

	data, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal products: %w", err)
	}

	if err := utils.WriteToFile(filename, data); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
