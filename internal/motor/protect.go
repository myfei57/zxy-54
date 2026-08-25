package motor

func (c *Controller) Protect(beltID string, temp float64) bool {
	if c.monitor == nil {
		return false
	}
	verdict := c.monitor.Filter().Evaluate(beltID, temp, c.monitor.Limit())
	return verdict.Over
}
