package store

import "minebelt/internal/belt"

type BeltSnapshot struct {
	ID      string
	Name    string
	Kind    string
	Section string
	Zone    string
	Status  int
	Speed   float64
	Load    float64
}

func (s *Store) SaveBeltSnapshots(rows []BeltSnapshot) error {
	return s.SaveJSON("belts.json", rows)
}

func (s *Store) LoadBeltSnapshots() ([]BeltSnapshot, error) {
	var rows []BeltSnapshot
	if err := s.LoadJSON("belts.json", &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func FromBelt(b *belt.Belt) BeltSnapshot {
	return BeltSnapshot{
		ID:      b.ID,
		Name:    b.Name,
		Kind:    b.Kind,
		Section: b.Section,
		Zone:    b.Zone,
		Status:  int(b.Status),
		Speed:   b.Speed,
		Load:    b.Load,
	}
}

func (s *BeltSnapshot) Apply(b *belt.Belt) {
	b.ID = s.ID
	b.Name = s.Name
	b.Kind = s.Kind
	b.Section = s.Section
	b.Zone = s.Zone
	b.Status = belt.Status(s.Status)
	b.Speed = s.Speed
	b.Load = s.Load
}
