package verifycase

import (
	"testing"

	"minebelt/internal/belt"
	"minebelt/internal/coal"
	"minebelt/internal/motor"
	"minebelt/internal/ns"
	"minebelt/internal/quota"
	"minebelt/internal/roller"
)

func TestMbBeltStartOrder(t *testing.T) {
	reg := ns.NewRegistry()
	line := belt.NewLine(reg)
	mon := roller.NewMonitor(90, 3)
	ctrl := motor.NewController(line, mon, nil)
	ledger := quota.NewLedger()
	feed := coal.NewFeederManager(ledger)
	b := line.AddBelt("b1", "main", "sec1")
	_ = ledger.AddQuota(b.ID, 1000, "shift")
	_ = feed.AddFeeder(b.ID, 100)
	ctrl.SetFeederPort(feed)
	feed.SetBeltReadyGate(ctrl)
	_ = line.SetStatus(b.ID, belt.Fault)
	_ = ctrl.StartOneKey(b.ID)
	if feed.IsRunning(b.ID) {
		t.Fatalf("feeder started before the belt was running")
	}
}
