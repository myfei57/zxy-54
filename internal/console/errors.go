package console

import "errors"

var (
	errNameOrSectionMissing = errors.New("name and section are required")
	errBeltMissing          = errors.New("belt not found")
	errRollerMissing        = errors.New("roller not found")
	errSensorMissing        = errors.New("gas sensor not found")
	errMotorMissing         = errors.New("motor not found")
	errQuotaMissing         = errors.New("quota not found")
	errStoreMissing         = errors.New("store is not configured")
)
