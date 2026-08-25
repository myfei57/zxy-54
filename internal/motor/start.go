package motor

import "minebelt/internal/belt"

func (c *Controller) Start(beltID string) error {
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
	return c.line.SetStatus(beltID, belt.Running)
}

func (c *Controller) StartOneKey(beltID string) error {
	if c.feeders != nil {
		if err := c.feeders.Start(beltID); err != nil {
			return err
		}
	}
	if err := c.Start(beltID); err != nil {
		return err
	}
	if c.feeders != nil {
		c.feeders.RecordOrder(beltID, "belt")
	}
	if err := c.ConfirmReady(beltID); err != nil {
		return err
	}
	if c.brakes != nil {
		if err := c.brakes.Release(beltID); err != nil {
			return err
		}
	}
	return nil
}
