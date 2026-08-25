package motor

func (c *Controller) TempVerdict(beltID string, temp float64) bool {
	base, ok := c.baseline[beltID]
	if !ok {
		base = 90
	}
	return temp > base
}

func (c *Controller) ReplaceMotor(beltID string, ratingKW, baseline float64) (*Motor, error) {
	m, ok := c.MotorForBelt(beltID)
	if !ok {
		return nil, errBeltNotFound
	}
	m.RatingKW = ratingKW
	m.Baseline = baseline
	if c.ratings != nil {
		_ = c.ratings.SaveRating(beltID, baseline)
	}
	return m, nil
}
