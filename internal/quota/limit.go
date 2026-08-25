package quota

func (l *Ledger) SetLimit(beltID string, limit float64) error {
	q, ok := l.quotas[beltID]
	if !ok {
		return quotaErrNotFound
	}
	q.Limit = limit
	return nil
}

func (l *Ledger) Remaining(beltID string) (float64, error) {
	q, ok := l.quotas[beltID]
	if !ok {
		return 0, quotaErrNotFound
	}
	rem := q.Limit - q.Used
	if rem < 0 {
		return 0, nil
	}
	return rem, nil
}
