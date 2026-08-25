package interlock

func (x *Interlock) Reset(beltID string) error {
	if x.latch == nil {
		return nil
	}
	return x.latch.ResetLatch(beltID)
}
