package codec

import (
	"encoding/json"
	"errors"

	"banned-products-system/internal/models"
)

var ErrInvalidType = errors.New("invalid type provided to codec")

// ProductCodec реализует goka.Codec для модели Product
type ProductCodec struct{}

// Encode сериализует Product в JSON
func (c *ProductCodec) Encode(value interface{}) ([]byte, error) {
	product, ok := value.(*models.Product)
	if !ok {
		return nil, ErrInvalidType
	}
	return json.Marshal(product)
}

// Decode десериализует JSON в Product
func (c *ProductCodec) Decode(data []byte) (interface{}, error) {
	var product models.Product
	err := json.Unmarshal(data, &product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}
