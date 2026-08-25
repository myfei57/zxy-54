package gas

import "fmt"

func (m *Manager) Calibrate(id string, threshold float64) error {
	s, ok := m.sensors[id]
	if !ok {
		return fmt.Errorf("sensor %s not found", id)
	}
	if threshold <= 0 {
		return fmt.Errorf("invalid threshold %f", threshold)
	}
	s.Threshold = threshold
	s.Version++
	if m.onCalibrate != nil {
		m.onCalibrate(id)
	}
	return nil
}

func (m *Manager) Threshold(id string) (float64, error) {
	s, ok := m.sensors[id]
	if !ok {
		return 0, fmt.Errorf("sensor %s not found", id)
	}
	return s.Threshold, nil
}
