package store

import "minebelt/internal/quota"

type QuotaSnapshot struct {
	BeltID string
	Limit  float64
	Used   float64
	Period string
}

func (s *Store) SaveQuotaSnapshots(rows []QuotaSnapshot) error {
	return s.SaveJSON("quotas.json", rows)
}

func (s *Store) LoadQuotaSnapshots() ([]QuotaSnapshot, error) {
	var rows []QuotaSnapshot
	if err := s.LoadJSON("quotas.json", &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func FromQuota(q *quota.Quota) QuotaSnapshot {
	return QuotaSnapshot{
		BeltID: q.BeltID,
		Limit:  q.Limit,
		Used:   q.Used,
		Period: q.Period,
	}
}

func (s *QuotaSnapshot) Apply(q *quota.Quota) {
	q.BeltID = s.BeltID
	q.Limit = s.Limit
	q.Used = s.Used
	q.Period = s.Period
}
