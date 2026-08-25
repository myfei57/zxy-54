package quota

func (l *Ledger) ResetPeriod(beltID string) error {
	q, ok := l.quotas[beltID]
	if !ok {
		return quotaErrNotFound
	}
	q.Used = 0
	return nil
}

func (l *Ledger) Rollover(beltID, period string) error {
	q, ok := l.quotas[beltID]
	if !ok {
		return quotaErrNotFound
	}
	q.Period = period
	q.Used = 0
	return nil
}

func (l *Ledger) Usage() []Quota {
	out := make([]Quota, 0, len(l.quotas))
	for _, q := range l.quotas {
		out = append(out, *q)
	}
	return out
}
