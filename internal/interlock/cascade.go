package interlock

func (x *Interlock) Cascade(beltIDs []string) []error {
	var errs []error
	for _, id := range beltIDs {
		if x.stopper == nil {
			continue
		}
		if err := x.stopper.Stop(id); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
