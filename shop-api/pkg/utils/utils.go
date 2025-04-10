package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// WriteToFile writes data to a file
func WriteToFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}

// ReadJSONFromFile reads a JSON file line by line
func ReadJSONFromFile(filename string) (<-chan []byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	// Create a channel to send JSON objects
	jsonChan := make(chan []byte)

	go func() {
		defer file.Close()
		defer close(jsonChan)

		decoder := json.NewDecoder(file)

		// Read opening bracket
		_, err := decoder.Token()
		if err != nil {
			fmt.Printf("Error reading opening token: %v\n", err)
			return
		}

		// Read array elements
		for decoder.More() {
			var product map[string]interface{}
			if err := decoder.Decode(&product); err != nil {
				fmt.Printf("Error decoding JSON: %v\n", err)
				continue
			}

			// Convert product back to JSON bytes
			productBytes, err := json.Marshal(product)
			if err != nil {
				fmt.Printf("Error encoding JSON: %v\n", err)
				continue
			}

			// Send to channel
			jsonChan <- productBytes

			// Small delay to simulate processing
			time.Sleep(100 * time.Millisecond)
		}

		// Read closing bracket
		_, err = decoder.Token()
		if err != nil {
			fmt.Printf("Error reading closing token: %v\n", err)
			return
		}
	}()

	return jsonChan, nil
}
