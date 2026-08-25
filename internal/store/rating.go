package store

func (s *Store) LoadRatings() (map[string]float64, error) {
	rows := make(map[string]float64)
	if err := s.LoadJSON("ratings.json", &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *Store) SaveRating(beltID string, baseline float64) error {
	rows := make(map[string]float64)
	_ = s.LoadJSON("ratings.json", &rows)
	rows[beltID] = baseline
	return s.SaveJSON("ratings.json", rows)
}
