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
	// 逆煤流启动：先启皮带，确认运行后再投给煤机，避免煤堆压在停带上。
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
	// 皮带已确认运行、制动已释放，方可投煤。
	if c.feeders != nil {
		if err := c.feeders.Start(beltID); err != nil {
			return err
		}
	}
	return nil
}
