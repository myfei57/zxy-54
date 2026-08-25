package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

var errFileNotFound = errors.New("file not found")

func (s *Store) path(name string) string {
	return filepath.Join(s.dir, name)
}

func (s *Store) SaveJSON(name string, value interface{}) error {
	if err := s.ensureDir(); err != nil {
		return err
	}
	target := s.path(name)
	tmp := target + ".tmp"
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, target)
}

func (s *Store) LoadJSON(name string, value interface{}) error {
	data, err := os.ReadFile(s.path(name))
	if err != nil {
		if os.IsNotExist(err) {
			return errFileNotFound
		}
		return err
	}
	if err := json.Unmarshal(data, value); err != nil {
		return err
	}
	return nil
}
