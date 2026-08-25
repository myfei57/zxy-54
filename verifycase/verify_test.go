package verifycase

import (
	"testing"

	"minebelt/internal/belt"
	"minebelt/internal/interlock"
	"minebelt/internal/motor"
	"minebelt/internal/ns"
	"minebelt/internal/roller"
)

func TestMbBatchStopCascade(t *testing.T) {
	reg := ns.NewRegistry()
	line := belt.NewLine(reg)
	mon := roller.NewMonitor(90, 3)
	ctrl := motor.NewController(line, mon, nil)
	b1 := line.AddBelt("b1", "main", "s1")
	b2 := line.AddBelt("b2", "main", "s2")
	b3 := line.AddBelt("b3", "main", "s3")
	_ = line.SetStatus(b1.ID, belt.Running)
	_ = line.SetStatus(b2.ID, belt.Fault)
	_ = line.SetStatus(b3.ID, belt.Running)
	xlock := interlock.New(nil, ctrl, nil, nil)
	_ = xlock.Cascade([]string{b1.ID, b2.ID, b3.ID})
	if line.IsRunning(b3.ID) {
		t.Fatalf("cascade stop halted on the first failure and left the trailing belt running")
	}
}
