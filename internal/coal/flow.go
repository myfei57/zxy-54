package coal

import "minebelt/internal/belt"

type Transit struct {
	beltID   string
	scale    *belt.Scale
	inst     float64
	total    float64
	peak     float64
	interval float64
}

func NewTransit(beltID string, scale *belt.Scale, interval float64) *Transit {
	if interval <= 0 {
		interval = 1
	}
	return &Transit{beltID: beltID, scale: scale, interval: interval}
}

func (t *Transit) Tick() {
	if t.scale.Pending() {
		t.inst = t.scale.Read()
	}
	t.total += t.inst * t.interval
	if t.inst > t.peak {
		t.peak = t.inst
	}
}

func (t *Transit) BeltID() string {
	return t.beltID
}

func (t *Transit) Inst() float64 {
	return t.inst
}

func (t *Transit) Total() float64 {
	return t.total
}

func (t *Transit) Peak() float64 {
	return t.peak
}

func (t *Transit) Average() float64 {
	if t.total == 0 {
		return 0
	}
	return t.total / t.interval
}

func (t *Transit) ResetTotal() {
	t.total = 0
	t.peak = 0
}
