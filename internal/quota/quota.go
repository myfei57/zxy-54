package quota

import "github.com/google/uuid"

type Quota struct {
	ID     string
	BeltID string
	Limit  float64
	Used   float64
	Period string
}

var ErrQuotaExhausted = quotaError("transport quota exhausted")

type Ledger struct {
	quotas map[string]*Quota
}

func NewLedger() *Ledger {
	return &Ledger{quotas: make(map[string]*Quota)}
}

func (l *Ledger) AddQuota(beltID string, limit float64, period string) *Quota {
	q := &Quota{ID: uuid.NewString(), BeltID: beltID, Limit: limit, Period: period}
	l.quotas[beltID] = q
	return q
}

func (l *Ledger) Quota(beltID string) (*Quota, bool) {
	q, ok := l.quotas[beltID]
	return q, ok
}

type quotaError string

func (e quotaError) Error() string {
	return string(e)
}
