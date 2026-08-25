package verifycase

import (
	"testing"

	"minebelt/internal/belt"
	"minebelt/internal/brake"
	"minebelt/internal/motor"
	"minebelt/internal/ns"
	"minebelt/internal/roller"
)

func TestMbBrakeReleaseOrder(t *testing.T) {
	reg := ns.NewRegistry()
	line := belt.NewLine(reg)
	mon := roller.NewMonitor(90, 3)
	ctrl := motor.NewController(line, mon, nil)
	b := line.AddBelt("b1", "main", "sec1")
	_ = line.SetStatus(b.ID, belt.Stopped)
	mgr := brake.NewManager()
	_ = mgr.AddBrake(b.ID)
	mgr.SetReadyPort(ctrl)
	_ = mgr.Release(b.ID)
	st, _ := mgr.Status(b.ID)
	if !st.Engaged {
		t.Fatalf("brake released before belt readiness was confirmed")
	}
}
