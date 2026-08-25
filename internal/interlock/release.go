package interlock

func (x *Interlock) Tag(beltID string) error {
	x.tagged[beltID] = true
	return nil
}

func (x *Interlock) ClearTag(beltID string) error {
	delete(x.tagged, beltID)
	return nil
}

func (x *Interlock) Release(beltID string) error {
	if x.tagged[beltID] {
		return errTagPresent
	}
	if x.power != nil {
		if err := x.power.IssueClearance(beltID); err != nil {
			return err
		}
		return x.power.Power(beltID)
	}
	return nil
}

var errTagPresent = interlockError("maintenance tag is still present")

type interlockError string

func (e interlockError) Error() string {
	return string(e)
}
