package motor

var (
	errBeltNotFound = motorError("belt not found")
	errNoLine       = motorError("belt line is not configured")
	errBeltFault    = motorError("belt is in fault state")
)

type motorError string

func (e motorError) Error() string {
	return string(e)
}
