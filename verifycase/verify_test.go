package verifycase

import (
	"testing"

	"minebelt/internal/gas"
	"minebelt/internal/interlock"
)

func TestMbGasInterlockThresholdFresh(t *testing.T) {
	gasm := gas.NewManager()
	s := gasm.AddSensor("face-a", 1.0)
	xlock := interlock.New(gasm, nil, nil, nil)
	gasm.SetCalibrateListener(func(string) { xlock.SyncThresholds() })
	_ = gasm.SetConcentration(s.ID, 0.8)
	_ = gasm.Calibrate(s.ID, 0.5)
	if !xlock.GasVerdict(s.ID) {
		t.Fatalf("gas interlock kept the old threshold after calibration")
	}
}
