package verifycase

import (
	"testing"

	"minebelt/internal/brake"
	"minebelt/internal/interlock"
)

func TestMbBrakeOverloadLatchReset(t *testing.T) {
	mgr := brake.NewManager()
	_ = mgr.AddBrake("b1")
	xlock := interlock.New(nil, nil, nil, mgr)
	_ = mgr.Trip("b1", "overcurrent")
	err := xlock.Reset("b1")
	if err == nil {
		t.Fatalf("overload latch reset while the overload cause was still active")
	}
	if !mgr.Latched("b1") {
		t.Fatalf("overload latch cleared while the cause remained active")
	}
}
