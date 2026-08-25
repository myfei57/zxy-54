package motor

import "minebelt/internal/belt"

var errBeltNotRunning = motorError("belt is not running")

func (c *Controller) Stop(beltID string) error {
	if c.line == nil {
		return errNoLine
	}
	b, ok := c.line.Belt(beltID)
	if !ok {
		return errBeltNotFound
	}
	if b.Status == belt.Fault {
		return errBeltFault
	}
	if b.Status == belt.Stopped {
		return errBeltNotRunning
	}
	if c.feeders != nil {
		_ = c.feeders.Stop(beltID)
	}
	if c.brakes != nil {
		_ = c.brakes.Engage(beltID)
	}
	return c.line.SetStatus(beltID, belt.Stopped)
}
