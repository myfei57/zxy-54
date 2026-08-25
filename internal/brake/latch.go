package brake

func (m *Manager) Latched(beltID string) bool {
	b, ok := m.brakes[beltID]
	if !ok {
		return false
	}
	return b.Overloaded
}

func (m *Manager) Cause(beltID string) string {
	b, ok := m.brakes[beltID]
	if !ok {
		return ""
	}
	return b.Cause
}

func (m *Manager) ResetLatch(beltID string) error {
	b, ok := m.brakes[beltID]
	if !ok {
		return errBrakeNotFound
	}
	if b.Cause == "" {
		return errCauseActive
	}
	b.Cause = ""
	b.Overloaded = false
	b.Engaged = false
	return nil
}

var errCauseActive = brakeError("overload cause is still active")
