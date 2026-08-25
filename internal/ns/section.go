package ns

type SectionKind int

const (
	KindMain SectionKind = iota
	KindExtended
	KindTransfer
)

func (k SectionKind) String() string {
	switch k {
	case KindMain:
		return "main"
	case KindExtended:
		return "extended"
	case KindTransfer:
		return "transfer"
	}
	return "unknown"
}

type Section struct {
	ID      string
	Name    string
	Kind    SectionKind
	Order   int
	BeltIDs []string
}
