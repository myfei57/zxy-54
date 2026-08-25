package interlock

func (x *Interlock) GasVerdict(sensorID string) bool {
	threshold, ok := x.thresholds[sensorID]
	if !ok {
		return false
	}
	conc, err := x.gas.Concentration(sensorID)
	if err != nil {
		return false
	}
	return conc > threshold
}
