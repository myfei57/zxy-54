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
	// 重载启动时必须先确认皮带就绪，再释放制动；否则制动松开后皮带会在
	// 负荷作用下倒溜，撞击尾部挡煤板。这里真正落下就绪标记，供制动释放前校验。
	if !c.line.IsReady(beltID) {
		return errNotReadyBelt
	}
	c.readySeq[beltID]++
	return nil
}
