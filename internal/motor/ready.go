package motor

var (
	errNotReadyBelt    = motorError("belt is not ready for start confirmation")
	errBrakeNotEngaged = motorError("brake is not engaged before start confirmation")
)

func (c *Controller) Ready(beltID string) bool {
	if c.line == nil {
		return false
	}
	return c.readySeq[beltID] > 0 && c.line.IsReady(beltID)
}

func (c *Controller) ReadyVersion(beltID string) int {
	return c.readySeq[beltID]
}

func (c *Controller) ConfirmReady(beltID string) error {
	if c.line == nil {
		return errNoLine
	}
	if _, ok := c.line.Belt(beltID); !ok {
		return errBeltNotFound
	}
	if !c.line.IsReady(beltID) {
		return errNotReadyBelt
	}
	if c.brakes != nil && !c.brakes.Engaged(beltID) {
		return errBrakeNotEngaged
	}
	c.readySeq[beltID]++
	return nil
}
