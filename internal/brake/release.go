package brake

type releaseRecord struct {
	beltID string
	stage  string
}

func (m *Manager) Release(beltID string) error {
	b, ok := m.brakes[beltID]
	if !ok {
		return errBrakeNotFound
	}
	b.Engaged = false
	m.recordRelease(beltID, "released")
	return nil
}

func (m *Manager) recordRelease(beltID, stage string) {
	m.releaseLog[beltID] = append(m.releaseLog[beltID], releaseRecord{beltID: beltID, stage: stage})
}

func (m *Manager) ReleaseLog(beltID string) []string {
	out := make([]string, 0, len(m.releaseLog[beltID]))
	for _, rec := range m.releaseLog[beltID] {
		out = append(out, rec.stage)
	}
	return out
}

func (m *Manager) ReleaseVersion(beltID string) int {
	return m.releaseSeq[beltID]
}

var errNotReady = brakeError("belt is not ready for brake release")
