package brake

import "fmt"

type Pressure struct {
	BeltID string
	Bar    float64
	Ready  bool
}

func (m *Manager) SetPressure(beltID string, bar float64) (Pressure, error) {
	b, ok := m.brakes[beltID]
	if !ok {
		return Pressure{}, errBrakeNotFound
	}
	if bar < 0 || bar > 120 {
		return Pressure{}, fmt.Errorf("invalid brake pressure %f", bar)
	}
	b.Pressure = bar
	return Pressure{BeltID: beltID, Bar: bar, Ready: bar >= 40}, nil
}

func (m *Manager) Pressure(beltID string) (Pressure, error) {
	b, ok := m.brakes[beltID]
	if !ok {
		return Pressure{}, errBrakeNotFound
	}
	return Pressure{BeltID: beltID, Bar: b.Pressure, Ready: b.Pressure >= 40}, nil
}
