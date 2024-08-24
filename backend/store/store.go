package store

import (
	"encoding/gob"
	"os"
	"sync"
)

const filename = ".save.bruh"

var mu = new(sync.Mutex)

// SaveBytes saves bytes to the save file.
func SaveBytes(b []byte) error {
	mu.Lock()
	defer mu.Unlock()

	// Create save file if it doesn't exist
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = gob.NewEncoder(file).Encode(b)
	if err != nil {
		return err
	}

	return nil
}

// ReadBytes reads bytes from the save file.
func ReadBytes() ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var contents []byte
	err = gob.NewDecoder(file).Decode(&contents)
	if err != nil {
		return nil, err
	}
	return contents, nil
}
