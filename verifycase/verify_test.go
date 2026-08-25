package verifycase

import (
	"testing"

	"minebelt/internal/belt"
	"minebelt/internal/coal"
	"minebelt/internal/ns"
)

func TestMbCoalFlowOrder(t *testing.T) {
	reg := ns.NewRegistry()
	line := belt.NewLine(reg)
	b := line.AddBelt("b1", "main", "sec1")
	sc, _ := line.Scale(b.ID)
	tr := coal.NewTransit(b.ID, sc, 1)
	sc.SetInstant(100)
	tr.Tick()
	sc.SetInstant(150)
	tr.Tick()
	if tr.Total() != 250 {
		t.Fatalf("flow total disagrees with the scale readings: got %v want 250", tr.Total())
	}
}
