package brake

type State struct {
	BeltID     string
	Engaged    bool
	Overloaded bool
	Cause      string
	Latch      bool
}

func (m *Manager) Status(beltID string) (State, error) {
	b, ok := m.brakes[beltID]
	if !ok {
		return State{}, errBrakeNotFound
	}
	return State{
		BeltID:     b.BeltID,
		Engaged:    b.Engaged,
		Overloaded: b.Overloaded,
		Cause:      b.Cause,
		Latch:      b.Overloaded,
	}, nil
}
