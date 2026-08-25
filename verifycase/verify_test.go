package verifycase

import (
	"testing"

	"minebelt/internal/belt"
	"minebelt/internal/interlock"
	"minebelt/internal/motor"
	"minebelt/internal/ns"
	"minebelt/internal/roller"
)

func TestMbClearanceOrder(t *testing.T) {
	reg := ns.NewRegistry()
	line := belt.NewLine(reg)
	mon := roller.NewMonitor(90, 3)
	ctrl := motor.NewController(line, mon, nil)
	b := line.AddBelt("b1", "main", "s1")
	xlock := interlock.New(nil, nil, ctrl, nil)
	_ = xlock.Tag(b.ID)
	_ = xlock.Release(b.ID)
	if ctrl.Powered(b.ID) {
		t.Fatalf("the belt powered up while the maintenance tag was still present")
	}
}
