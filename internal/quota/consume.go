package quota

func (l *Ledger) Consume(beltID string, tons float64) error {
	q, ok := l.quotas[beltID]
	if !ok {
		return quotaErrNotFound
	}
	q.Used += tons
	return nil
}

func (l *Ledger) Exhausted(beltID string) bool {
	q, ok := l.quotas[beltID]
	if !ok {
		return false
	}
	return q.Used >= q.Limit
}

var quotaErrNotFound = quotaError("quota not found")
