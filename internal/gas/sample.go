package gas

import "fmt"

func (m *Manager) SetConcentration(id string, value float64) error {
	s, ok := m.sensors[id]
	if !ok {
		return fmt.Errorf("sensor %s not found", id)
	}
	s.Concentration = value
	return nil
}

func (m *Manager) Concentration(id string) (float64, error) {
	s, ok := m.sensors[id]
	if !ok {
		return 0, fmt.Errorf("sensor %s not found", id)
	}
	return s.Concentration, nil
}
