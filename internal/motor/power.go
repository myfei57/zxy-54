package motor

import "minebelt/internal/belt"

var errNoClearance = motorError("belt has no maintenance clearance")

func (c *Controller) Power(beltID string) error {
	if c.line == nil {
		return errNoLine
	}
	if _, ok := c.line.Belt(beltID); !ok {
		return errBeltNotFound
	}
	if !c.VerifyClearance(beltID) {
		return errNoClearance
	}
	c.ClearClearance(beltID)
	c.powered[beltID] = true
	if m, ok := c.MotorForBelt(beltID); ok {
		m.Powered = true
	}
	c.recordClearance(beltID, "powered")
	return c.line.SetStatus(beltID, belt.Running)
}

func (c *Controller) Powered(beltID string) bool {
	return c.powered[beltID]
}

func (c *Controller) IssueClearance(beltID string) error {
	if _, ok := c.line.Belt(beltID); !ok {
		return errBeltNotFound
	}
	c.clearances[beltID] = beltID
	c.recordClearance(beltID, "issued")
	return nil
}

func (c *Controller) VerifyClearance(beltID string) bool {
	return c.clearances[beltID] == beltID
}

func (c *Controller) ClearClearance(beltID string) {
	delete(c.clearances, beltID)
	c.recordClearance(beltID, "consumed")
}

func (c *Controller) ClearanceLog(beltID string) []string {
	return append([]string(nil), c.clearanceLog[beltID]...)
}

func (c *Controller) recordClearance(beltID, stage string) {
	c.clearanceLog[beltID] = append(c.clearanceLog[beltID], stage)
}
