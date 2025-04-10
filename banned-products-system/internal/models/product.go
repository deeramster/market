package models

import "time"

// Product представляет модель товара
type Product struct {
	Name           string    `json:"name"`
	Price          float64   `json:"price"`
	Category       string    `json:"category"`
	Brand          string    `json:"brand"`
	Stock          Stock     `json:"stock"`
	SKU            string    `json:"sku"`
	Tags           []string  `json:"tags"`
	Images         []Image   `json:"images"`
	Specifications Specs     `json:"specifications"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Index          string    `json:"index"`
	StoreID        string    `json:"store_id"`
}

// Stock содержит информацию о наличии товара
type Stock struct {
	Available int `json:"available"`
	Reserved  int `json:"reserved"`
}

// Image содержит информацию об изображении товара
type Image struct {
	URL string `json:"url"`
	Alt string `json:"alt"`
}

// Specs содержит технические характеристики товара
type Specs struct {
	Weight          string `json:"weight"`
	Dimensions      string `json:"dimensions"`
	BatteryLife     string `json:"battery_life"`
	WaterResistance string `json:"water_resistance"`
}

// BannedProduct представляет запрещенный товар
type BannedProduct struct {
	SKU      string    `json:"sku"`
	Reason   string    `json:"reason"`
	BannedAt time.Time `json:"banned_at"`
}

// BannedProductsList представляет список запрещенных товаров
type BannedProductsList struct {
	Products []BannedProduct `json:"products"`
}
