package gas

import "github.com/google/uuid"

type Sensor struct {
	ID            string
	FaceID        string
	Concentration float64
	Threshold     float64
	Version       int
}

type Manager struct {
	sensors     map[string]*Sensor
	onCalibrate func(sensorID string)
	alarms      []alarmRecord
}

func NewManager() *Manager {
	return &Manager{sensors: make(map[string]*Sensor)}
}

func (m *Manager) AddSensor(face string, threshold float64) *Sensor {
	s := &Sensor{ID: uuid.NewString(), FaceID: face, Threshold: threshold}
	m.sensors[s.ID] = s
	return s
}

func (m *Manager) SensorByFace(face string) (*Sensor, bool) {
	for _, s := range m.sensors {
		if s.FaceID == face {
			return s, true
		}
	}
	return nil, false
}

func (m *Manager) SetCalibrateListener(fn func(sensorID string)) {
	m.onCalibrate = fn
}

func (m *Manager) Sensors() []string {
	ids := make([]string, 0, len(m.sensors))
	for id := range m.sensors {
		ids = append(ids, id)
	}
	return ids
}
