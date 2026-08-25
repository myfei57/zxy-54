package store

import "os"

type Store struct {
	dir string
}

func NewStore(dir string) (*Store, error) {
	s := &Store{dir: dir}
	if err := s.ensureDir(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) Dir() string {
	return s.dir
}

func (s *Store) ensureDir() error {
	return os.MkdirAll(s.dir, 0o755)
}
