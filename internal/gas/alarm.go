package gas

import "time"

type Alarm struct {
	FaceID        string
	Concentration float64
	Threshold     float64
	At            time.Time
}

type alarmRecord struct {
	FaceID        string
	Concentration float64
	Threshold     float64
	At            time.Time
}

func (m *Manager) RecordAlarm(face string, concentration, threshold float64) Alarm {
	rec := alarmRecord{
		FaceID:        face,
		Concentration: concentration,
		Threshold:     threshold,
		At:            time.Now().UTC(),
	}
	m.alarms = append(m.alarms, rec)
	if len(m.alarms) > 100 {
		m.alarms = m.alarms[len(m.alarms)-100:]
	}
	return Alarm(rec)
}

func (m *Manager) Alarms(n int) []Alarm {
	if n < 1 {
		n = 20
	}
	out := make([]Alarm, 0, n)
	for i := len(m.alarms) - 1; i >= 0 && len(out) < n; i-- {
		out = append(out, Alarm(m.alarms[i]))
	}
	return out
}

func (m *Manager) AlarmCount() int {
	return len(m.alarms)
}
