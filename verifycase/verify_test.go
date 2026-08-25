package verifycase

import (
	"testing"

	"minebelt/internal/belt"
	"minebelt/internal/motor"
	"minebelt/internal/ns"
	"minebelt/internal/roller"
	"minebelt/internal/store"
)

func TestMbMotorTempBaselineFresh(t *testing.T) {
	dir := t.TempDir()
	st, err := store.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	reg := ns.NewRegistry()
	line := belt.NewLine(reg)
	mon := roller.NewMonitor(90, 3)
	ctrl := motor.NewController(line, mon, st)
	b := line.AddBelt("b1", "main", "s1")
	_ = ctrl.AddMotor(b.ID, 250, 90)
	_, _ = ctrl.ReplaceMotor(b.ID, 400, 110)
	if ctrl.TempVerdict(b.ID, 100) {
		t.Fatalf("the new motor tripped on a normal temperature using the old baseline")
	}
}
