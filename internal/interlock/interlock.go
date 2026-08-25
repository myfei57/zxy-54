package interlock

type ThresholdSource interface {
	Threshold(sensorID string) (float64, error)
	Concentration(sensorID string) (float64, error)
	Sensors() []string
}

type Stopper interface {
	Stop(beltID string) error
}

type PowerPort interface {
	Power(beltID string) error
	Powered(beltID string) bool
	IssueClearance(beltID string) error
	VerifyClearance(beltID string) bool
}

type LatchPort interface {
	ResetLatch(beltID string) error
	Latched(beltID string) bool
	Cause(beltID string) string
}

type Interlock struct {
	gas        ThresholdSource
	stopper    Stopper
	power      PowerPort
	latch      LatchPort
	thresholds map[string]float64
	tagged     map[string]bool
}

func New(gas ThresholdSource, stopper Stopper, power PowerPort, latch LatchPort) *Interlock {
	x := &Interlock{
		gas:        gas,
		stopper:    stopper,
		power:      power,
		latch:      latch,
		thresholds: make(map[string]float64),
		tagged:     make(map[string]bool),
	}
	x.SyncThresholds()
	return x
}

func (x *Interlock) SyncThresholds() {
	if x.gas == nil {
		return
	}
	for _, id := range x.gas.Sensors() {
		if v, err := x.gas.Threshold(id); err == nil {
			x.thresholds[id] = v
		}
	}
}
