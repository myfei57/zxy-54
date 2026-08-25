package coal

func (m *FeederManager) RecordProduction(beltID string, tons float64) error {
	if m.ledger == nil {
		return nil
	}
	return m.ledger.Consume(beltID, tons)
}
