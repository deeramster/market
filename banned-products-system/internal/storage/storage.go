package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"banned-products-system/internal/models"
)

// Storage представляет хранилище запрещенных товаров
type Storage struct {
	filePath string
	list     models.BannedProductsList
	mu       sync.RWMutex
}

// New создает новый экземпляр хранилища
func New(filePath string) (*Storage, error) {
	s := &Storage{
		filePath: filePath,
		list: models.BannedProductsList{
			Products: []models.BannedProduct{},
		},
	}

	// Загружаем список запрещенных товаров из файла
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return s, nil
}

// load загружает список запрещенных товаров из файла
func (s *Storage) load() error {
	file, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Если файл не существует, используем пустой список
			s.list = models.BannedProductsList{Products: []models.BannedProduct{}}
			return nil
		}
		return err
	}

	return json.Unmarshal(file, &s.list)
}

// save сохраняет список запрещенных товаров в файл
func (s *Storage) save() error {
	data, err := json.MarshalIndent(s.list, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}

// Add добавляет товар в список запрещенных
func (s *Storage) Add(sku, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверка дубликатов
	for _, product := range s.list.Products {
		if product.SKU == sku {
			return fmt.Errorf("товар с SKU %s уже в списке запрещенных", sku)
		}
	}

	s.list.Products = append(s.list.Products, models.BannedProduct{
		SKU:      sku,
		Reason:   reason,
		BannedAt: time.Now(),
	})

	return s.save()
}

// Remove удаляет товар из списка запрещенных
func (s *Storage) Remove(sku string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, product := range s.list.Products {
		if product.SKU == sku {
			// Удаляем элемент из списка
			s.list.Products = append(s.list.Products[:i], s.list.Products[i+1:]...)
			return s.save()
		}
	}

	return fmt.Errorf("товар с SKU %s не найден в списке запрещенных", sku)
}

// Update обновляет причину запрета товара
func (s *Storage) Update(sku, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, product := range s.list.Products {
		if product.SKU == sku {
			s.list.Products[i].Reason = reason
			s.list.Products[i].BannedAt = time.Now()
			return s.save()
		}
	}

	return fmt.Errorf("товар с SKU %s не найден в списке запрещенных", sku)
}

// List возвращает список запрещенных товаров
func (s *Storage) List() []models.BannedProduct {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Создаем копию списка, чтобы избежать проблем с конкурентным доступом
	products := make([]models.BannedProduct, len(s.list.Products))
	copy(products, s.list.Products)

	return products
}

// IsBanned проверяет, находится ли товар в списке запрещенных
func (s *Storage) IsBanned(sku string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, product := range s.list.Products {
		if product.SKU == sku {
			return true
		}
	}

	return false
}
