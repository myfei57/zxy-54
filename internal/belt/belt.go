package belt

type Status int

const (
	Stopped Status = iota
	Running
	Fault
)

func (s Status) String() string {
	switch s {
	case Running:
		return "running"
	case Fault:
		return "fault"
	}
	return "stopped"
}

type Belt struct {
	ID      string
	Name    string
	Kind    string
	Section string
	Zone    string
	Status  Status
	Speed   float64
	Load    float64
}
