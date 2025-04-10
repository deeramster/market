package models

import "time"

// Price represents the product price
type Price struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// Stock represents the product stock information
type Stock struct {
	Available int `json:"available"`
	Reserved  int `json:"reserved"`
}

// Image represents a product image
type Image struct {
	URL string `json:"url"`
	Alt string `json:"alt"`
}

// Specifications represents product specifications
type Specifications struct {
	Weight          string `json:"weight"`
	Dimensions      string `json:"dimensions"`
	BatteryLife     string `json:"battery_life"`
	WaterResistance string `json:"water_resistance"`
}

// Product represents the product information
type Product struct {
	ProductID      string         `json:"product_id"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Price          Price          `json:"price"`
	Category       string         `json:"category"`
	Brand          string         `json:"brand"`
	Stock          Stock          `json:"stock"`
	SKU            string         `json:"sku"`
	Tags           []string       `json:"tags"`
	Images         []Image        `json:"images"`
	Specifications Specifications `json:"specifications"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	Index          string         `json:"index"`
	StoreID        string         `json:"store_id"`
}
