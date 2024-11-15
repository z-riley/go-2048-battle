package store

import (
	"encoding/gob"
	"os"
	"sync"
)

// Store can store bytes to the disk.
type Store struct {
	mu       *sync.Mutex
	filename string
}

// NewStore constructs a new store under a specified filename. If the file already
// exists, its current contents are used.
func NewStore(filename string) *Store {
	return &Store{
		mu:       new(sync.Mutex),
		filename: filename,
	}
}

// SaveBytes saves bytes to the store.
func (s *Store) SaveBytes(b []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create save file if it doesn't exist
	file, err := os.Create(s.filename)
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

// ReadBytes reads bytes from the store.
func (s *Store) ReadBytes() ([]byte, error) {
	file, err := os.Open(s.filename)
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
