package belt

import "fmt"

func (l *Line) SetSpeed(id string, speed float64) error {
	b, ok := l.belts[id]
	if !ok {
		return fmt.Errorf("belt %s not found", id)
	}
	if speed < 0 {
		return fmt.Errorf("invalid speed %f", speed)
	}
	b.Speed = speed
	return nil
}

func (l *Line) SetLoad(id string, load float64) error {
	b, ok := l.belts[id]
	if !ok {
		return fmt.Errorf("belt %s not found", id)
	}
	if load < 0 {
		return fmt.Errorf("invalid load %f", load)
	}
	b.Load = load
	return nil
}

func (l *Line) RunningCount() int {
	count := 0
	for _, b := range l.belts {
		if b.Status == Running {
			count++
		}
	}
	return count
}

func (l *Line) FaultCount() int {
	count := 0
	for _, b := range l.belts {
		if b.Status == Fault {
			count++
		}
	}
	return count
}

type TransferPoint struct {
	ID      string
	Name    string
	From    string
	To      string
	Blocked bool
}

type TransferState struct {
	points map[string]*TransferPoint
}

func NewTransferState() *TransferState {
	return &TransferState{points: make(map[string]*TransferPoint)}
}

func (t *TransferState) AddPoint(id, name, from, to string) *TransferPoint {
	p := &TransferPoint{ID: id, Name: name, From: from, To: to}
	t.points[id] = p
	return p
}

func (t *TransferState) Point(id string) (*TransferPoint, bool) {
	p, ok := t.points[id]
	return p, ok
}

func (t *TransferState) SetBlocked(id string, blocked bool) error {
	p, ok := t.points[id]
	if !ok {
		return fmt.Errorf("transfer point %s not found", id)
	}
	p.Blocked = blocked
	return nil
}

func (t *TransferState) All() []*TransferPoint {
	out := make([]*TransferPoint, 0, len(t.points))
	for _, p := range t.points {
		out = append(out, p)
	}
	return out
}
