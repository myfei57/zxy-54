package interlock

// Cascade issues a stop to every belt in the chain. A cascade must bring down
// the whole line: even if one belt errors (for example, the head belt that is
// already in Fault), the remaining downstream belts still have to stop, or the
// coal flow piles up at the transfer point. Per-belt errors are therefore
// accumulated and reported, never used to short-circuit the loop.
func (x *Interlock) Cascade(beltIDs []string) []error {
	var errs []error
	for _, id := range beltIDs {
		if x.stopper == nil {
			break
		}
		if err := x.stopper.Stop(id); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
