package brake

import "github.com/google/uuid"

type Brake struct {
	ID         string
	BeltID     string
	Engaged    bool
	Overloaded bool
	Cause      string
	Pressure   float64
}

type ReadyPort interface {
	Ready(beltID string) bool
	ReadyVersion(beltID string) int
}

type Manager struct {
	brakes     map[string]*Brake
	ready      ReadyPort
	releaseLog map[string][]releaseRecord
	releaseSeq map[string]int
}

func NewManager() *Manager {
	return &Manager{
		brakes:     make(map[string]*Brake),
		releaseLog: make(map[string][]releaseRecord),
		releaseSeq: make(map[string]int),
	}
}

func (m *Manager) SetReadyPort(p ReadyPort) {
	m.ready = p
}

func (m *Manager) AddBrake(beltID string) *Brake {
	b := &Brake{ID: uuid.NewString(), BeltID: beltID, Engaged: true}
	m.brakes[beltID] = b
	return b
}

func (m *Manager) Brake(beltID string) (*Brake, bool) {
	b, ok := m.brakes[beltID]
	return b, ok
}

func (m *Manager) Trip(beltID, cause string) error {
	b, ok := m.brakes[beltID]
	if !ok {
		return errBrakeNotFound
	}
	b.Cause = cause
	b.Overloaded = true
	b.Engaged = true
	return nil
}

func (m *Manager) Engage(beltID string) error {
	b, ok := m.brakes[beltID]
	if !ok {
		return errBrakeNotFound
	}
	b.Engaged = true
	return nil
}

func (m *Manager) Engaged(beltID string) bool {
	b, ok := m.brakes[beltID]
	return ok && b.Engaged
}

var errBrakeNotFound = brakeError("brake not found")

type brakeError string

func (e brakeError) Error() string {
	return string(e)
}
