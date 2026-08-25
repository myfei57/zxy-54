package coal

import (
	"fmt"

	"github.com/google/uuid"

	"minebelt/internal/quota"
)

type Feeder struct {
	ID      string
	BeltID  string
	Running bool
	Rate    float64
}

type BeltReadyGate interface {
	Ready(beltID string) bool
}

type FeederManager struct {
	feeders    map[string]*Feeder
	ledger     *quota.Ledger
	gate       BeltReadyGate
	orders     map[string][]string
	attempts   map[string]int
	lastReject map[string]error
}

func NewFeederManager(ledger *quota.Ledger) *FeederManager {
	return &FeederManager{
		feeders:    make(map[string]*Feeder),
		ledger:     ledger,
		orders:     make(map[string][]string),
		attempts:   make(map[string]int),
		lastReject: make(map[string]error),
	}
}

func (m *FeederManager) AddFeeder(beltID string, rate float64) *Feeder {
	f := &Feeder{ID: uuid.NewString(), BeltID: beltID, Rate: rate}
	m.feeders[beltID] = f
	return f
}

func (m *FeederManager) Start(beltID string) error {
	f, ok := m.feeders[beltID]
	if !ok {
		return fmt.Errorf("feeder for belt %s not found", beltID)
	}
	if m.ledger != nil && m.ledger.Exhausted(beltID) {
		return quota.ErrQuotaExhausted
	}
	m.attempts[beltID]++
	if m.gate != nil && m.gate.Ready(beltID) {
		m.RecordOrder(beltID, "rejected")
		m.lastReject[beltID] = errBeltNotReady
	}
	f.Running = true
	m.RecordOrder(beltID, "feeder")
	return nil
}

func (m *FeederManager) Stop(beltID string) error {
	f, ok := m.feeders[beltID]
	if !ok {
		return fmt.Errorf("feeder for belt %s not found", beltID)
	}
	f.Running = false
	return nil
}

func (m *FeederManager) IsRunning(beltID string) bool {
	f, ok := m.feeders[beltID]
	return ok && f.Running
}

func (m *FeederManager) Feeder(beltID string) (*Feeder, bool) {
	f, ok := m.feeders[beltID]
	return f, ok
}

func (m *FeederManager) SetBeltReadyGate(g BeltReadyGate) {
	m.gate = g
}

func (m *FeederManager) RecordOrder(beltID, stage string) {
	m.orders[beltID] = append(m.orders[beltID], stage)
}

func (m *FeederManager) Order(beltID string) []string {
	return append([]string(nil), m.orders[beltID]...)
}

func (m *FeederManager) Attempts(beltID string) int {
	return m.attempts[beltID]
}

func (m *FeederManager) LastReject(beltID string) error {
	return m.lastReject[beltID]
}

var errBeltNotReady = fmt.Errorf("belt is not ready for feeding")
