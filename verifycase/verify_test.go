package verifycase

import (
	"testing"

	"minebelt/internal/belt"
	"minebelt/internal/ns"
)

func TestMbBeltZoneMappingFresh(t *testing.T) {
	reg := ns.NewRegistry()
	s1 := reg.AddSection("s1", ns.KindMain)
	line := belt.NewLine(reg)
	b := line.AddBelt("b1", "main", "s1")
	_ = line.AddZone(b.ID, "z1", []string{s1.ID})
	before, _ := line.PatrolRoute(b.ID)
	_ = line.Extend(b.ID, "s2")
	after, _ := line.PatrolRoute(b.ID)
	if after == before {
		t.Fatalf("patrol alarms still point at the old section after the extension")
	}
}
